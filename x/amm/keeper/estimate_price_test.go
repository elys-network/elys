package keeper_test

import (
	"cosmossdk.io/math"
	elystypes "github.com/elys-network/elys/types"
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
				price, _ := suite.app.AmmKeeper.GetTokenPrice(suite.ctx, ptypes.BaseCurrency, ptypes.BaseCurrency)
				suite.Require().Equal(elystypes.ZeroDec34().String(), price.String())
			},
		},
		{
			"get token price with oracle price",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices()
			},
			func() {
				price, _ := suite.app.AmmKeeper.GetTokenPrice(suite.ctx, ptypes.BaseCurrency, ptypes.BaseCurrency)
				suite.Require().Equal(elystypes.OneDec34().String(), price.String())
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
				price, _ := suite.app.AmmKeeper.GetTokenPrice(suite.ctx, ptypes.BaseCurrency, ptypes.BaseCurrency)
				suite.Require().Equal(elystypes.ZeroDec34().String(), price.String())
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
				value := suite.app.AmmKeeper.CalculateUSDValue(suite.ctx, ptypes.BaseCurrency, math.NewInt(1000))
				suite.Require().Equal(value.String(), elystypes.NewDec34FromInt64(1000).String())
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
				suite.Require().Equal(value.String(), elystypes.ZeroDec34().String())
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
