package types

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (p *Pool) CalculateWeightFees(ctx sdk.Context, oracleKeeper OracleKeeper,
	accountedAssets []PoolAsset,
	finalAssetsPool []PoolAsset, tokenDenom string, params Params, weightBreakingFeePerpetualFactor sdkmath.LegacyDec,
) (sdkmath.LegacyDec, sdkmath.LegacyDec) {

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, finalAssetsPool)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// target weight
	targetWeightIn := GetDenomNormalizedWeight(p.PoolAssets, tokenDenom)
	targetWeightOut := sdkmath.LegacyOneDec().Sub(targetWeightIn)

	// weight breaking fee as in Plasma pool
	finalWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, finalAssetsPool, tokenDenom)
	finalWeightOut := sdkmath.LegacyOneDec().Sub(finalWeightIn)

	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, accountedAssets, tokenDenom)
	initialWeightOut := sdkmath.LegacyOneDec().Sub(initialWeightIn)
	weightBreakingFee := GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut, distanceDiff, params)
	// weightBreakingFeePerpetualFactor is 1 if not send by perpetual
	weightBreakingFee = weightBreakingFee.Mul(weightBreakingFeePerpetualFactor)
	// weight recovery reward = weight breaking fee * weight breaking fee portion
	weightRecoveryReward := weightBreakingFee.Mul(params.WeightBreakingFeePortion)

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus := weightBreakingFee.Neg()

	if distanceDiff.IsNegative() {
		weightBreakingFee = sdkmath.LegacyZeroDec()
		weightBalanceBonus = sdkmath.LegacyZeroDec()

		// set weight breaking fee to zero if bonus is applied
		if initialWeightDistance.GT(params.ThresholdWeightDifference) {
			weightBalanceBonus = weightRecoveryReward
		}

		if initialWeightDistance.GT(params.ThresholdWeightDifferenceSwapFee) {
			swapFee = sdkmath.LegacyZeroDec()
		}
	}
	if initialWeightDistance.GT(params.ThresholdWeightDifference) && distanceDiff.IsNegative() {
		weightBalanceBonus = weightRecoveryReward
		// set weight breaking fee to zero if bonus is applied
		weightBreakingFee = sdkmath.LegacyZeroDec()
	}

	return weightBalanceBonus, weightBreakingFee
}
