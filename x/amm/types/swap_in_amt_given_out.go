package types

import (
	"fmt"

	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p Pool) CalcGivenOutSlippage(
	ctx sdk.Context,
	oracleKeeper OracleKeeper,
	snapshot *Pool,
	tokensOut sdk.Coins,
	tokenInDenom string,
	accPoolKeeper AccountedPoolKeeper,
) (sdkmath.LegacyDec, error) {
	balancerInCoin, _, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, sdkmath.LegacyZeroDec(), accPoolKeeper)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	// in amount is calculated in this formula
	oracleInAmount := sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)
	balancerIn := sdkmath.LegacyNewDecFromInt(balancerInCoin.Amount)
	balancerSlippage := balancerIn.Sub(oracleInAmount)
	if balancerSlippage.IsNegative() {
		return sdkmath.LegacyZeroDec(), nil
	}
	return balancerSlippage, nil
}

// SwapInAmtGivenOut is a mutative method for CalcOutAmtGivenIn, which includes the actual swap.
// weightBreakingFeePerpetualFactor should be 1 if perpetual is not the one calling this function
func (p *Pool) SwapInAmtGivenOut(
	ctx sdk.Context, oracleKeeper OracleKeeper, snapshot *Pool,
	tokensOut sdk.Coins, tokenInDenom string, swapFee sdkmath.LegacyDec, accPoolKeeper AccountedPoolKeeper, weightBreakingFeePerpetualFactor math.LegacyDec) (
	tokenIn sdk.Coin, slippage, slippageAmount sdkmath.LegacyDec, weightBalanceBonus sdkmath.LegacyDec, err error,
) {
	// early return with balancer swap if normal amm pool
	if !p.PoolParams.UseOracle {
		balancerInCoin, slippage, err := p.CalcInAmtGivenOut(ctx, oracleKeeper, snapshot, tokensOut, tokenInDenom, swapFee, accPoolKeeper)
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		err = p.applySwap(ctx, sdk.Coins{balancerInCoin}, tokensOut, swapFee, sdkmath.LegacyZeroDec())
		if err != nil {
			return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
		}

		return balancerInCoin, slippage, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), nil
	}

	tokenOut, poolAssetOut, poolAssetIn, err := p.parsePoolAssets(tokensOut, tokenInDenom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}

	// ensure token prices for in/out tokens set properly
	inTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenInDenom)
	if inTokenPrice.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for inToken not set: %s", poolAssetIn.Token.Denom)
	}
	outTokenPrice := oracleKeeper.GetAssetPriceFromDenom(ctx, tokenOut.Denom)
	if outTokenPrice.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), fmt.Errorf("price for outToken not set: %s", poolAssetOut.Token.Denom)
	}

	accountedAssets := p.GetAccountedBalance(ctx, accPoolKeeper, p.PoolAssets)
	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)

	// in amount is calculated in this formula
	// balancer slippage amount = Max(oracleOutAmount-balancerOutAmount, 0)
	// resizedAmount = tokenIn / externalLiquidityRatio
	// actualSlippageAmount = balancer slippage(resizedAmount)
	oracleInAmount := sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Mul(outTokenPrice).Quo(inTokenPrice)

	externalLiquidityRatio, err := p.GetAssetExternalLiquidityRatio(tokenOut.Denom)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	// Ensure externalLiquidityRatio is not zero to avoid division by zero
	if externalLiquidityRatio.IsZero() {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrAmountTooLow
	}

	resizedAmount := sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Quo(externalLiquidityRatio).RoundInt()
	slippageAmount, err = p.CalcGivenOutSlippage(
		ctx,
		oracleKeeper,
		snapshot,
		sdk.Coins{sdk.NewCoin(tokenOut.Denom, resizedAmount)},
		tokenInDenom,
		accPoolKeeper,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	inAmountAfterSlippage := oracleInAmount.Add(slippageAmount)
	slippage = slippageAmount.Quo(oracleInAmount)

	// calculate weight distance difference to calculate bonus/cut on the operation
	newAssetPools, err := p.NewPoolAssetsAfterSwap(ctx,
		sdk.Coins{sdk.NewCoin(tokenInDenom, inAmountAfterSlippage.TruncateInt())},
		tokensOut, accountedAssets,
	)
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, newAssetPools)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// target weight
	// Asset weight remains same in new pool assets as in original pool assets
	targetWeightIn := GetDenomNormalizedWeight(newAssetPools, tokenInDenom)
	targetWeightOut := GetDenomNormalizedWeight(newAssetPools, tokenOut.Denom)

	// weight breaking fee as in Plasma pool
	weightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenInDenom)
	weightOut := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, newAssetPools, tokenOut.Denom)
	weightBreakingFee := GetWeightBreakingFee(weightIn, weightOut, targetWeightIn, targetWeightOut, p.PoolParams, distanceDiff)

	// weightBreakingFeePerpetualFactor is 1 if not send by perpetual
	weightBreakingFee = weightBreakingFee.Mul(weightBreakingFeePerpetualFactor)
	// weight recovery reward = weight breaking fee * weight recovery fee portion
	weightRecoveryReward := weightBreakingFee.Mul(p.PoolParams.WeightRecoveryFeePortion)

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus = weightBreakingFee.Neg()

	// If swap is improving weight, set weight breaking fee to zero
	if distanceDiff.IsNegative() {
		weightBreakingFee = sdkmath.LegacyZeroDec()
		weightBalanceBonus = sdkmath.LegacyZeroDec()

		// set weight breaking fee to zero if bonus is applied
		if initialWeightDistance.GT(p.PoolParams.ThresholdWeightDifference) {
			weightBalanceBonus = weightRecoveryReward
		}
	}

	if swapFee.GTE(sdkmath.LegacyOneDec()) {
		return sdk.Coin{}, sdkmath.LegacyOneDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), ErrTooMuchSwapFee
	}

	tokenAmountInInt := inAmountAfterSlippage.
		Mul(sdkmath.LegacyOneDec().Add(weightBreakingFee)).
		Quo(sdkmath.LegacyOneDec().Sub(swapFee)).
		TruncateInt()
	tokenIn = sdk.NewCoin(tokenInDenom, tokenAmountInInt)
	err = p.applySwap(ctx, sdk.Coins{tokenIn}, tokensOut, swapFee, sdkmath.LegacyZeroDec())
	if err != nil {
		return sdk.Coin{}, sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), sdkmath.LegacyZeroDec(), err
	}
	return tokenIn, slippage, slippageAmount, weightBalanceBonus, nil
}
