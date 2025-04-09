package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func PoolAddress() sdk.AccAddress {
	return authtypes.NewModuleAddress(ModuleName)
}

func (p *AmmPool) AddLiabilities(coin sdk.Coin) {
	p.TotalLiabilities = p.TotalLiabilities.Add(coin)
}

func (p *AmmPool) SubLiabilities(coin sdk.Coin) {
	p.TotalLiabilities = p.TotalLiabilities.Sub(coin)
}

func (p Pool) GetBigDecMaxWithdrawRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaxWithdrawRatio)
}

func (p Pool) GetBigDecMaxLeverageRatio() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.MaxLeverageRatio)
}

func (p Pool) GetBigDecInterestRate() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.InterestRate)
}

func (p Pool) GetBigDecInterestRateMax() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.InterestRateMax)
}

func (p Pool) GetBigDecInterestRateMin() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.InterestRateMin)
}

func (p Pool) GetBigDecInterestRateIncrease() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.InterestRateIncrease)
}

func (p Pool) GetBigDecInterestRateDecrease() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.InterestRateDecrease)
}

func (p Pool) GetBigDecHealthGainFactor() osmomath.BigDec {
	return osmomath.BigDecFromDec(p.HealthGainFactor)
}

func (p Pool) GetBigDecTotalValue() osmomath.BigDec {
	return osmomath.BigDecFromSDKInt(p.TotalValue)
}
