package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (p Pool) CalcInAmtGivenOut(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdk.Dec, accountedPool AccountedPoolKeeper) (
	tokenIn sdk.Coin, slippage sdk.Dec, err error,
) {
	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	outWeight := sdk.NewDecFromInt(poolAssetOut.Weight)
	inWeight := sdk.NewDecFromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetOut, poolAssetIn, err := snapshot.parsePoolAssets(tokensOut, tokenInDenom)
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
		}
		// accounted pool balance
		acountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
		if acountedPoolAssetOutAmt.IsPositive() {
			poolAssetIn.Token.Amount = acountedPoolAssetOutAmt
		}

		// accounted pool balance
		acountedPoolAssetOutAmt = accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
		if acountedPoolAssetOutAmt.IsPositive() {
			poolAssetOut.Token.Amount = acountedPoolAssetOutAmt
		}
		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, sdk.ZeroDec(), err
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
		return sdk.Coin{}, sdk.ZeroDec(), err
	}
	tokenAmountIn = tokenAmountIn.Neg()

	rate, err := p.GetTokenARate(ctx, oracle, snapshot, tokenInDenom, tokenOut.Denom, accountedPool)
	if err != nil {
		return sdk.Coin{}, sdk.ZeroDec(), err
	}

	amountInWithoutSlippage := sdk.NewDecFromInt(tokenOut.Amount).Quo(rate)
	if tokenAmountIn.IsZero() {
		return sdk.Coin{}, sdk.ZeroDec(), ErrAmountTooLow
	}
	slippage = sdk.OneDec().Sub(amountInWithoutSlippage.Quo(tokenAmountIn))

	// Ensure (1 - swapfee) is not zero to avoid division by zero
	if swapFee.GTE(sdk.OneDec()) {
		return sdk.Coin{}, sdk.ZeroDec(), ErrTooMuchSwapFee
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
		return sdk.Coin{}, sdk.ZeroDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "token amount must be positive")
	}
	return sdk.NewCoin(tokenInDenom, tokenInAmt), slippage, nil
}
