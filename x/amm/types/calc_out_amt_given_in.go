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
	accountedPool AccountedPoolKeeper,
) (sdk.Coin, osmomath.BigDec, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), err
	}

	tokenAmountInAfterFee := osmomath.BigDecFromSDKInt(tokenIn.Amount).Mul(osmomath.OneBigDec().Sub(swapFee))
	poolTokenInBalance := osmomath.BigDecFromSDKInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = osmomath.BigDecFromSDKInt(acountedPoolAssetInAmt)
	}

	poolTokenOutBalance := osmomath.BigDecFromSDKInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	accountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if accountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = osmomath.BigDecFromSDKInt(accountedPoolAssetOutAmt)
	}

	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	outWeight := osmomath.BigDecFromSDKInt(poolAssetOut.Weight)
	inWeight := osmomath.BigDecFromSDKInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetIn, poolAssetOut, err := snapshot.parsePoolAssets(tokensIn, tokenOutDenom)
		if err != nil {
			return sdk.Coin{}, osmomath.ZeroBigDec(), err
		}
		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, osmomath.ZeroBigDec(), err
		}
		inWeight = oracleWeights[0].Weight
		outWeight = oracleWeights[1].Weight
	}

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
		return sdk.Coin{}, osmomath.ZeroBigDec(), err
	}

	if tokenAmountOut.IsZero() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), ErrTokenOutAmountZero
	}

	rate, err := p.GetTokenARate(ctx, oracle, tokenIn.Denom, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, osmomath.ZeroBigDec(), err
	}

	amountOutWithoutSlippage := tokenAmountInAfterFee.Mul(rate)

	// check if amountOutWithoutSlippage is zero to avoid division by zero
	if amountOutWithoutSlippage.IsZero() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "amount out without slippage must be positive")
	}

	slippage := osmomath.OneBigDec().Sub(tokenAmountOut.Quo(amountOutWithoutSlippage))

	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.Dec().TruncateInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, osmomath.ZeroBigDec(), ErrTokenOutAmountZero
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), slippage, nil
}
