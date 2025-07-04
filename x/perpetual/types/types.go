package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func GetPositionFromString(s string) Position {
	switch s {
	case "long":
		return Position_LONG
	case "short":
		return Position_SHORT
	default:
		return Position_UNSPECIFIED
	}
}

func (i InterestBlock) GetBigDecInterestRate() osmomath.BigDec {
	return osmomath.BigDecFromDec(i.InterestRate)
}

func (f FundingRateBlock) GetBigDecFundingRateLong() osmomath.BigDec {
	return osmomath.BigDecFromDec(f.FundingRateLong)
}

func (f FundingRateBlock) GetBigDecFundingRateShort() osmomath.BigDec {
	return osmomath.BigDecFromDec(f.FundingRateShort)
}

func (f FundingRateBlock) GetBigDecFundingShareLong() osmomath.BigDec {
	return osmomath.BigDecFromDec(f.FundingShareLong)
}

func (f FundingRateBlock) GetBigDecFundingShareShort() osmomath.BigDec {
	return osmomath.BigDecFromDec(f.FundingShareShort)
}

func NewPerpetualFees(perpFees, slippageFees, weightBreakingFees, takerFees sdk.Coins) PerpetualFees {
	return PerpetualFees{
		PerpFees:           perpFees,
		SlippageFees:       slippageFees,
		WeightBreakingFees: weightBreakingFees,
		TakerFees:          takerFees,
	}
}

func NewPerpetualFeesWithEmptyCoins() PerpetualFees {
	return PerpetualFees{
		PerpFees:           sdk.Coins{},
		SlippageFees:       sdk.Coins{},
		WeightBreakingFees: sdk.Coins{},
		TakerFees:          sdk.Coins{},
	}
}

func (p PerpetualFees) Add(other PerpetualFees) PerpetualFees {
	return NewPerpetualFees(
		sdk.Coins(p.PerpFees).Sort().Add(sdk.Coins(other.PerpFees).Sort()...),
		sdk.Coins(p.SlippageFees).Sort().Add(sdk.Coins(other.SlippageFees).Sort()...),
		sdk.Coins(p.WeightBreakingFees).Add(sdk.Coins(other.WeightBreakingFees).Sort()...),
		sdk.Coins(p.TakerFees).Add(sdk.Coins(other.TakerFees).Sort()...),
	)
}
