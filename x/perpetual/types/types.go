package types

import (
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
