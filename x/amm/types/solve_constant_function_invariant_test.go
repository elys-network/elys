package types_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/x/amm/types"
)

func (suite *TestSuite) TestCalculateTokenARate() {
	tokenARate := types.CalculateTokenARate(sdkmath.LegacyNewDec(10), sdkmath.LegacyNewDec(3), sdkmath.LegacyNewDec(100), sdkmath.LegacyNewDec(6))
	suite.Require().Equal(tokenARate.String(), sdkmath.LegacyNewDec(5).String())

	tokenARate = types.CalculateTokenARate(sdkmath.LegacyNewDec(0), sdkmath.LegacyNewDec(3), sdkmath.LegacyNewDec(10), sdkmath.LegacyNewDec(6))
	suite.Require().Equal(tokenARate.String(), sdkmath.LegacyNewDec(0).String())

	tokenARate = types.CalculateTokenARate(sdkmath.LegacyNewDec(10), sdkmath.LegacyNewDec(3), sdkmath.LegacyNewDec(100), sdkmath.LegacyNewDec(0))
	suite.Require().Equal(tokenARate.String(), sdkmath.LegacyNewDec(0).String())
}
