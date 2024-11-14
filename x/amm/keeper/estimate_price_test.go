package keeper_test

import (
	"cosmossdk.io/math"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *AmmKeeperTestSuite) TestEstimatePrice() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"get token price without oracle price",
			func() {
				suite.ResetSuite()
			},
			func() {
				price := suite.app.AmmKeeper.GetTokenPrice(suite.ctx, ptypes.BaseCurrency, ptypes.BaseCurrency)
				suite.Require().Equal(math.LegacyZeroDec(), price)
			},
		},
		{
			"get token price with oracle price",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				price := suite.app.AmmKeeper.GetTokenPrice(suite.ctx, ptypes.BaseCurrency, ptypes.BaseCurrency)
				suite.Require().Equal(math.LegacyMustNewDecFromStr("0.000001000000000000"), price)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			tc.postValidateFunction()
		})
	}
}
