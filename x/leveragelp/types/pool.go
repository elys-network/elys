package types

import (
	sdkmath "cosmossdk.io/math"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func NewPool(poolId uint64, maxLeverage, maxLeveragelpRatio, adlTriggerRatio sdkmath.LegacyDec) Pool {
	return Pool{
		AmmPoolId:          poolId,
		Health:             sdkmath.LegacyOneDec(),
		LeveragedLpAmount:  sdkmath.ZeroInt(),
		LeverageMax:        maxLeverage,
		MaxLeveragelpRatio: maxLeveragelpRatio,
		AdlTriggerRatio:    adlTriggerRatio,
	}
}

func (p *Pool) UpdateAssetLeveragedAmount(denom string, amount sdkmath.Int, isIncrease bool) {
	newAssetLevAmounts := make([]*AssetLeverageAmount, 0)
	for _, asset := range p.AssetLeverageAmounts {
		if asset.Denom == denom {
			if isIncrease {
				asset.LeveragedAmount = asset.LeveragedAmount.Add(amount)
			} else {
				asset.LeveragedAmount = asset.LeveragedAmount.Sub(amount)
			}
		}
		newAssetLevAmounts = append(newAssetLevAmounts, asset)
	}
	p.AssetLeverageAmounts = newAssetLevAmounts
}

func (p Pool) GetBigDecLeveragedLpAmount() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(p.LeveragedLpAmount)
}

func (p Pool) GetBigDecMaxLeveragelpRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaxLeveragelpRatio)
}
