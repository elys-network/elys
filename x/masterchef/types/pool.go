package types

import "github.com/osmosis-labs/osmosis/osmomath"

func (p PoolInfo) GetBigDecMultiplier() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.Multiplier)
}

func (p PoolRewardsAccum) GetBigDecDexReward() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.DexReward)
}

func (p PoolRewardsAccum) GetBigDecGasReward() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.GasReward)
}

func (p PoolRewardsAccum) GetBigDecEdenReward() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.EdenReward)
}
