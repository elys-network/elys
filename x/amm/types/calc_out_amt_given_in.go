package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// CalcOutAmtGivenIn calculates tokens to be swapped out given the provided
// amount and fee deducted, using solveConstantFunctionInvariant.
func (p Pool) CalcOutAmtGivenIn(
	ctx sdk.Context,
	oracle OracleKeeper,
	snapshot *Pool,
	tokensIn sdk.Coins,
	tokenOutDenom string,
	swapFee sdkmath.LegacyDec,
	accountedPool AccountedPoolKeeper,
) (sdk.Coin, sdkmath.LegacyDec, error) {
	tokenIn, poolAssetIn, poolAssetOut, err := p.parsePoolAssets(tokensIn, tokenOutDenom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
	}

	tokenAmountInAfterFee := sdkmath.LegacyNewDecFromInt(tokenIn.Amount).Mul(sdkmath.LegacyOneDec().Sub(swapFee))
	poolTokenInBalance := sdkmath.LegacyNewDecFromInt(poolAssetIn.Token.Amount)
	// accounted pool balance
	acountedPoolAssetInAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetIn.Token.Denom)
	if acountedPoolAssetInAmt.IsPositive() {
		poolTokenInBalance = sdkmath.LegacyNewDecFromInt(acountedPoolAssetInAmt)
	}

	poolTokenOutBalance := sdkmath.LegacyNewDecFromInt(poolAssetOut.Token.Amount)
	// accounted pool balance
	accountedPoolAssetOutAmt := accountedPool.GetAccountedBalance(ctx, p.PoolId, poolAssetOut.Token.Denom)
	if accountedPoolAssetOutAmt.IsPositive() {
		poolTokenOutBalance = sdkmath.LegacyNewDecFromInt(accountedPoolAssetOutAmt)
	}

	poolPostSwapInBalance := poolTokenInBalance.Add(tokenAmountInAfterFee)

	outWeight := sdkmath.LegacyNewDecFromInt(poolAssetOut.Weight)
	inWeight := sdkmath.LegacyNewDecFromInt(poolAssetIn.Weight)
	if p.PoolParams.UseOracle {
		_, poolAssetIn, poolAssetOut, err := snapshot.parsePoolAssets(tokensIn, tokenOutDenom)
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
		}
		oracleWeights, err := GetOraclePoolNormalizedWeights(ctx, p.PoolId, oracle, []PoolAsset{poolAssetIn, poolAssetOut})
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
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
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
	}
	// REV: Origin could be here, improve error handling
	// This is the origin
	if tokenAmountOut.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), ErrTokenOutAmountZero
	}

	rate, err := p.GetTokenARate(ctx, oracle, snapshot, tokenIn.Denom, tokenOutDenom, accountedPool)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), err
	}

	amountOutWithoutSlippage := tokenAmountInAfterFee.Mul(rate)

	// check if amountOutWithoutSlippage is zero to avoid division by zero
	if amountOutWithoutSlippage.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), errorsmod.Wrapf(ErrInvalidMathApprox, "amount out without slippage must be positive")
	}

	slippage := sdkmath.LegacyOneDec().Sub(tokenAmountOut.Quo(amountOutWithoutSlippage))

	// We ignore the decimal component, as we round down the token amount out.
	tokenAmountOutInt := tokenAmountOut.TruncateInt()
	if !tokenAmountOutInt.IsPositive() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), ErrTokenOutAmountZero
	}

	return sdk.NewCoin(tokenOutDenom, tokenAmountOutInt), slippage, nil
}
