package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (p *Pool) CalculateWeightFees(ctx sdk.Context, oracleKeeper OracleKeeper,
	accountedAssets []PoolAsset,
	finalAssetsPool []PoolAsset, tokenInDenom string, params Params, weightBreakingFeePerpetualFactor osmomath.BigDec,
) (osmomath.BigDec, osmomath.BigDec, bool) {
	swapFee := true

	initialWeightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, accountedAssets)
	weightDistance := p.WeightDistanceFromTarget(ctx, oracleKeeper, finalAssetsPool)
	distanceDiff := weightDistance.Sub(initialWeightDistance)

	// target weight
	targetWeightIn := GetDenomNormalizedWeight(p.PoolAssets, tokenInDenom)
	targetWeightOut := osmomath.OneBigDec().Sub(targetWeightIn)

	// weight breaking fee as in Plasma pool
	finalWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, finalAssetsPool, tokenInDenom)
	finalWeightOut := osmomath.OneBigDec().Sub(finalWeightIn)

	initialWeightIn := GetDenomOracleAssetWeight(ctx, p.PoolId, oracleKeeper, accountedAssets, tokenInDenom)
	initialWeightOut := osmomath.OneBigDec().Sub(initialWeightIn)
	weightBreakingFee := GetWeightBreakingFee(finalWeightIn, finalWeightOut, targetWeightIn, targetWeightOut, initialWeightIn, initialWeightOut, distanceDiff, params)
	// weightBreakingFeePerpetualFactor is 1 if not send by perpetual
	weightBreakingFee = weightBreakingFee.Mul(weightBreakingFeePerpetualFactor)
	// weight recovery reward = weight breaking fee * weight breaking fee portion
	weightRecoveryReward := weightBreakingFee.Mul(params.GetBigDecWeightBreakingFeePortion())

	// bonus is valid when distance is lower than original distance and when threshold weight reached
	weightBalanceBonus := weightBreakingFee.Neg()

	if distanceDiff.IsNegative() {
		weightBreakingFee = osmomath.ZeroBigDec()
		weightBalanceBonus = osmomath.ZeroBigDec()

		// set weight breaking fee to zero if bonus is applied
		if initialWeightDistance.GT(params.GetBigDecThresholdWeightDifference()) {
			weightBalanceBonus = weightRecoveryReward
		}

		if initialWeightDistance.GT(params.GetBigDecThresholdWeightDifferenceSwapFee()) {
			swapFee = false
		}
	} else {
		// Weight getting worst but threshold is not reached so fees should not be charged
		if weightDistance.LT(params.GetBigDecThresholdWeightDifference()) {
			weightBreakingFee = osmomath.ZeroBigDec()
			weightBalanceBonus = osmomath.ZeroBigDec()
		}
	}

	return weightBalanceBonus, weightBreakingFee, swapFee
}
