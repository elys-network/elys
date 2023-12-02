package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (p Pool) CalcInAmtGivenOut(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec, accountedPool AccountedPoolKeeper) (
	tokenIn sdk.Coin, err error,
) {
	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, err
	}

	outWeight := sdk.NewDecFromInt(poolAssetOut.Weight)
	inWeight := sdk.NewDecFromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetOut, poolAssetIn, err := snapshot.parsePoolAssets(tokensOut, tokenInDenom)
		oracleWeights, err := OraclePoolNormalizedWeights(ctx, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, err
		}
		inWeight = oracleWeights[0].Weight
		outWeight = oracleWeights[1].Weight
	}

	// delta balanceOut is positive(tokens inside the pool decreases)
	poolTokenOutBalance := sdk.NewDecFromInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	acountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if acountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = sdk.NewDecFromInt(acountedPoolAssetOutAmt)
	}

	poolTokenInBalance := sdk.NewDecFromInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = sdk.NewDecFromInt(acountedPoolAssetInAmt)
	}

	poolPostSwapOutBalance := poolTokenOutBalance.Sub(sdk.NewDecFromInt(tokenOut.Amount))
	// (x_0)(y_0) = (x_0 + in)(y_0 - out)
	tokenAmountIn, err := solveConstantFunctionInvariant(
		poolTokenOutBalance, poolPostSwapOutBalance,
		outWeight,
		poolTokenInBalance,
		inWeight,
	)
	if err != nil {
		return sdk.Coin{}, err
	}
	tokenAmountIn = tokenAmountIn.Neg()

	// Ensure (1 - swapfee) is not zero to avoid division by zero
	if sdk.OneDec().Sub(swapFee).IsZero() {
		return sdk.Coin{}, ErrAmountTooLow
	}

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInBeforeFee := tokenAmountIn.Quo(sdk.OneDec().Sub(swapFee))

	// We round up tokenInAmt, as this is whats charged for the swap, for the precise amount out.
	// Otherwise, the pool would under-charge by this rounding error.
	tokenInAmt := tokenAmountInBeforeFee.Ceil().TruncateInt()

	if !tokenInAmt.IsPositive() {
		return sdk.Coin{}, sdkerrors.Wrapf(ErrInvalidMathApprox, "token amount must be positive")
	}
	return sdk.NewCoin(tokenInDenom, tokenInAmt), nil
}
