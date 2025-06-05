package types

import "github.com/osmosis-labs/osmosis/osmomath"

func (u UserRewardInfo) GetBigDecRewardDebt() osmomath.BigDec {
	return osmomath.BigDecFromDec(u.RewardDebt)
}

func (u UserRewardInfo) GetBigDecRewardPending() osmomath.BigDec {
	return osmomath.BigDecFromDec(u.RewardPending)
}

func (p PoolRewardInfo) GetBigDecPoolAccRewardPerShare() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.PoolAccRewardPerShare)
}
