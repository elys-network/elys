package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
)

func (p *Pool) CalculateWeightFees(ctx sdk.Context, oracleKeeper OracleKeeper,
	accountedAssets []PoolAsset,
	finalAssetsPool []PoolAsset, tokenDenom string, params Params, weightBreakingFeePerpetualFactor sdkmath.LegacyDec,
) (elystypes.Dec34, elystypes.Dec34, bool) {
	swapFee := true

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, finalAssetsPool)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// target weight
	targetWeightIn := GetDenomNormalizedWeight(p.PoolAssets, tokenDenom)
	targetWeightOut := elystypes.OneDec34().Sub(targetWeightIn)

	// weight breaking fee as in Plasma pool
	finalWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, finalAssetsPool, tokenDenom)
	finalWeightOut := elystypes.OneDec34().Sub(finalWeightIn)

	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, accountedAssets, tokenDenom)
	initialWeightOut := elystypes.OneDec34().Sub(initialWeightIn)
	weightBreakingFee := GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut, distanceDiff, params)
	// weightBreakingFeePerpetualFactor is 1 if not send by perpetual
	weightBreakingFee = weightBreakingFee.MulLegacyDec(weightBreakingFeePerpetualFactor)
	// weight recovery reward = weight breaking fee * weight breaking fee portion
	weightRecoveryReward := weightBreakingFee.MulLegacyDec(params.WeightBreakingFeePortion)

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus := weightBreakingFee.Neg()

	if distanceDiff.IsNegative() {
		weightBreakingFee = elystypes.ZeroDec34()
		weightBalanceBonus = elystypes.ZeroDec34()

		// set weight breaking fee to zero if bonus is applied
		if initialWeightDistance.GT(elystypes.NewDec34FromLegacyDec(params.ThresholdWeightDifference)) {
			weightBalanceBonus = weightRecoveryReward
		}

		if initialWeightDistance.GT(elystypes.NewDec34FromLegacyDec(params.ThresholdWeightDifferenceSwapFee)) {
			swapFee = false
		}
	} else {
		// Weight getting worst but threshold is not reached so fees should not be charged
		if initialWeightDistance.LT(elystypes.NewDec34FromLegacyDec(params.ThresholdWeightDifference)) {
			weightBreakingFee = elystypes.ZeroDec34()
			weightBalanceBonus = elystypes.ZeroDec34()
		}
	}

	return weightBalanceBonus, weightBreakingFee, swapFee
}
