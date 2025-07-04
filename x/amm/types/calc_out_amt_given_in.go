package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcOutAmtGivenIn calculates tokens to be swapped out given the provided
// amount and fee deducted, using solveConstantFunctionInvariant.
func (p Pool) CalcOutAmtGivenIn(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot SnapshotPool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee osmomath.BigDec,
) (sdk.Coin, osmomath.BigDec, osmomath.BigDec, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	outWeight := osmomath.BigDecFromSDKInt(poolAssetOut.Weight)
	inWeight := osmomath.BigDecFromSDKInt(poolAssetIn.Weight)

	if p.PoolParams.UseOracle {
		_, poolAssetIn, poolAssetOut, err = snapshot.parsePoolAssets(tokensIn, tokenOutDenom)
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

	tokenAmountInAfterFee := osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(osmomath.OneBigDec().Sub(swapFee))

	poolTokenInBalance := osmomath.BigDecFromSDKInt(poolAssetIn.Token.Amount)
	poolTokenOutBalance := osmomath.BigDecFromSDKInt(poolAssetOut.Token.Amount)

	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	// deduct swapfee on the tokensIn
	// delta balanceOut is positive(tokens inside the pool decreases)
	tokenAmountOut, err := solveConstantFunctionInvariant(
		poolTokenInBalance,
		poolPostSwapInBalance,
		inWeight,
		poolTokenOutBalance,
		outWeight,
	)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	if tokenAmountOut.IsZero() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ErrTokenOutAmountZero
	}

	rate, err := p.GetTokenARate(ctx, oracle, tokenIn.Denom, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	amountOutWithoutSlippage := tokenAmountInAfterFee.Mul(rate)

	// check if amountOutWithoutSlippage is zero to avoid division by zero
	if amountOutWithoutSlippage.IsZero() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "amount out without slippage must be positive")
	}

	slippage := osmomath.OneBigDec().Sub(tokenAmountOut.Quo(amountOutWithoutSlippage))

	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.Dec().TruncateInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), ErrTokenOutAmountZero
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), slippage, tokenAmountOut, nil
}
