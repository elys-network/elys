package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/elys-network/elys/x/oracle/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryPrice() {
	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(1),
			Source:    "elys",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(2),
			Source:    "band",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(3),
			Source:    "band",
			Timestamp: 110000,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	resp, err := suite.app.OracleKeeper.Price(suite.ctx, &types.QueryGetPriceRequest{
		Asset: "BTC",
	})
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Price, prices[2])
	resp, err = suite.app.OracleKeeper.Price(suite.ctx, &types.QueryGetPriceRequest{
		Asset:  "BTC",
		Source: "elys",
	})
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Price, prices[0])
	resp, err = suite.app.OracleKeeper.Price(suite.ctx, &types.QueryGetPriceRequest{
		Asset:     "BTC",
		Source:    "elys",
		Timestamp: 100000,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Price, prices[0])
	resp, err = suite.app.OracleKeeper.Price(suite.ctx, &types.QueryGetPriceRequest{
		Asset:     "BTC",
		Source:    "elys",
		Timestamp: 11000,
	})
	suite.Require().Error(err)
}
