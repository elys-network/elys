package types_test

import (
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (suite *TestSuite) TestCalculateTokenARate() {
	tokenARate := types.CalculateTokenARate(osmomath.NewBigDec(10), osmomath.NewBigDec(3), osmomath.NewBigDec(100), osmomath.NewBigDec(6))
	suite.Require().Equal(tokenARate.String(), osmomath.NewBigDec(5).String())

	tokenARate = types.CalculateTokenARate(osmomath.NewBigDec(0), osmomath.NewBigDec(3), osmomath.NewBigDec(10), osmomath.NewBigDec(6))
	suite.Require().Equal(tokenARate.String(), osmomath.NewBigDec(0).String())

	tokenARate = types.CalculateTokenARate(osmomath.NewBigDec(10), osmomath.NewBigDec(3), osmomath.NewBigDec(100), osmomath.NewBigDec(0))
	suite.Require().Equal(tokenARate.String(), osmomath.NewBigDec(0).String())
}
