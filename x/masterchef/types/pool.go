package types

import "github.com/osmosis-labs/osmosis/osmomath"

func (p PoolInfo) GetBigDecMultiplier() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.Multiplier)
}

func (u UserRewardInfo) GetBigDecRewardDebt() osmomath.BigDec {
	return osmomath.BigDecFromDec(u.RewardDebt)
}

func (u UserRewardInfo) GetBigDecRewardPending() osmomath.BigDec {
	return osmomath.BigDecFromDec(u.RewardPending)
}

func (p PoolRewardInfo) GetBigDecPoolAccRewardPerShare() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.PoolAccRewardPerShare)
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
