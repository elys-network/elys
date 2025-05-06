package types

import (
	"math"

	"github.com/osmosis-labs/osmosis/osmomath"
)

const UsdcPoolId = math.MaxInt16

func (i InterestBlock) GetBigDecInterestRate() osmomath.BigDec {
	return osmomath.BigDecFromDec(i.InterestRate)
}
