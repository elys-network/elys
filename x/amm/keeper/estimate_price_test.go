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
		{
			"Asset Info Not found for tokenInDenom",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				suite.app.OracleKeeper.RemoveAssetInfo(suite.ctx, ptypes.BaseCurrency)
				price := suite.app.AmmKeeper.GetTokenPrice(suite.ctx, ptypes.BaseCurrency, ptypes.BaseCurrency)
				suite.Require().Equal(math.LegacyMustNewDecFromStr("0.000000000000000000"), price)
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

func (suite *AmmKeeperTestSuite) TestCalculateUSDValue() {
	testCases := []struct {
		name                 string
		prerequisiteFunction func()
		postValidateFunction func()
	}{
		{
			"Success: get token value at oracle price",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				value := suite.app.AmmKeeper.CalculateUSDValue(suite.ctx, ptypes.ATOM, math.NewInt(1000))
				suite.Require().Equal(value, math.LegacyMustNewDecFromStr("0.001"))
			},
		},
		{
			"Calculate Usd value for asset not found in AssetProfile",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				value := suite.app.AmmKeeper.CalculateUSDValue(suite.ctx, "dummy", math.NewInt(1000))
				suite.Require().Equal(value.String(), math.LegacyZeroDec().String())
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
