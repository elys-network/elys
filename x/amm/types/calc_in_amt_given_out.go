package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcInAmtGivenOut calculates token to be provided, fee added,
// given the swapped out amount, using solveConstantFunctionInvariant.
func (p Pool) CalcInAmtGivenOut(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot SnapshotPool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee osmomath.BigDec) (
	tokenIn sdk.Coin, slippage osmomath.BigDec, tokenInDec osmomath.BigDec, err error,
) {
	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	outWeight := osmomath.BigDecFromSDKInt(poolAssetOut.Weight)
	inWeight := osmomath.BigDecFromSDKInt(poolAssetIn.Weight)

	if p.PoolParams.UseOracle {
		_, poolAssetOut, poolAssetIn, err = snapshot.parsePoolAssets(tokensOut, tokenInDenom)
		if err != nil {
			return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}
		inWeight = oracleWeights[0].Weight
		outWeight = oracleWeights[1].Weight
	}

	// delta balanceOut is positive(tokens inside the pool decreases)
	poolTokenOutBalance := osmomath.BigDecFromSDKInt(poolAssetOut.Token.Amount)
	poolTokenInBalance := osmomath.BigDecFromSDKInt(poolAssetIn.Token.Amount)

	poolPostSwapOutBalance := poolTokenOutBalance.Sub(osmomath.BigDecFromSDKInt(tokenOut.Amount))
	// (x_0)(y_0) = (x_0 + in)(y_0 - out)
	tokenAmountIn, err := solveConstantFunctionInvariant(
		poolTokenOutBalance, poolPostSwapOutBalance,
		outWeight,
		poolTokenInBalance,
		inWeight,
	)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	tokenAmountIn = tokenAmountIn.Neg()

	rate, err := p.GetTokenARate(ctx, oracle, tokenInDenom, tokenOut.Denom)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	amountInWithoutSlippage := osmomath.BigDecFromSDKInt(tokenOut.Amount).Quo(rate)
	if tokenAmountIn.IsZero() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ErrAmountTooLow
	}
	slippage = osmomath.OneBigDec().Sub(tokenAmountIn.Quo(amountInWithoutSlippage))

	// Ensure (1 - swapfee) is not zero to avoid division by zero
	if swapFee.GTE(osmomath.OneBigDec()) {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ErrTooMuchSwapFee
	}

	// We deduct a swap fee on the input asset. The swap happens by following the invariant curve on the input * (1 - swap fee)
	// and then the swap fee is added to the pool.
	// Thus in order to give X amount out, we solve the invariant for the invariant input. However invariant input = (1 - swapfee) * trade input.
	// Therefore we divide by (1 - swapfee) here
	tokenAmountInBeforeFee := tokenAmountIn.Quo(osmomath.OneBigDec().Sub(swapFee))

	// We round up tokenInAmt, as this is whats charged for the swap, for the precise amount out.
	// Otherwise, the pool would under-charge by this rounding error.
	tokenInAmt := tokenAmountInBeforeFee.Ceil().Dec().TruncateInt()

	if !tokenInAmt.IsPositive() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "token amount must be positive")
	}
	return sdk.NewCoin(tokenInDenom, tokenInAmt), slippage, tokenAmountInBeforeFee, nil
}
