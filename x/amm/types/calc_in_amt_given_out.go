package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (p Pool) CalcInAmtGivenOut(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdkmath.LegacyDec, accountedPool AccountedPoolKeeper) (
	tokenIn sdk.Coin, slippage elystypes.Dec34, err error,
) {
	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), err
	}

	outWeight := elystypes.NewDec34FromInt(poolAssetOut.Weight)
	inWeight := elystypes.NewDec34FromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetOut, poolAssetIn, err := snapshot.parsePoolAssets(tokensOut, tokenInDenom)
		if err != nil {
			return sdk.Coin{}, elystypes.ZeroDec34(), err
		}
		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, elystypes.ZeroDec34(), err
		}
		inWeight = oracleWeights[0].Weight
		outWeight = oracleWeights[1].Weight
	}

	// delta balanceOut is positive(tokens inside the pool decreases)
	poolTokenOutBalance := elystypes.NewDec34FromInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	acountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if acountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = elystypes.NewDec34FromInt(acountedPoolAssetOutAmt)
	}

	poolTokenInBalance := elystypes.NewDec34FromInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = elystypes.NewDec34FromInt(acountedPoolAssetInAmt)
	}

	poolPostSwapOutBalance := poolTokenOutBalance.SubInt(tokenOut.Amount)
	// (x_0)(y_0) = (x_0 + in)(y_0 - out)
	tokenAmountIn, err := solveConstantFunctionInvariant(
		poolTokenOutBalance, poolPostSwapOutBalance,
		outWeight,
		poolTokenInBalance,
		inWeight,
	)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), err
	}
	tokenAmountIn = tokenAmountIn.Neg()

	rate, err := p.GetTokenARate(ctx, oracle, snapshot, tokenInDenom, tokenOut.Denom, accountedPool)
	if err != nil {
		return sdk.Coin{}, elystypes.ZeroDec34(), err
	}

	amountInWithoutSlippage := elystypes.NewDec34FromInt(tokenOut.Amount).Quo(rate)
	if tokenAmountIn.IsZero() {
		return sdk.Coin{}, elystypes.ZeroDec34(), ErrAmountTooLow
	}
	slippage = elystypes.OneDec34().Sub(tokenAmountIn.Quo(amountInWithoutSlippage))

	// Ensure (1 - swapfee) is not zero to avoid division by zero
	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, elystypes.ZeroDec34(), ErrTooMuchSwapFee
	}

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInBeforeFee := tokenAmountIn.Quo(elystypes.OneDec34().SubLegacyDec(swapFee))

	// print tokenAmountInBeforeFee
	fmt.Println("tokenAmountInBeforeFee", tokenAmountInBeforeFee.String())

	// We round up tokenInAmt, as this is whats charged for the swap, for the precise amount out.
	// Otherwise, the pool would under-charge by this rounding error.
	tokenInAmt := tokenAmountInBeforeFee.ToInt()

	// print tokenInAmt
	fmt.Println("tokenInAmt", tokenInAmt.String())

	if !tokenInAmt.IsPositive() {
		return sdk.Coin{}, elystypes.ZeroDec34(), errorsmod.Wrapf(ErrInvalidMathApprox, "token amount must be positive")
	}
	return sdk.NewCoin(tokenInDenom, tokenInAmt), slippage, nil
}
