package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

func (suite *KeeperTestSuite) TestGRPCQueryPrice() {
	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(1),
			Source:    "binance",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(2),
			Source:    "band",
			Timestamp: 100000,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(3),
			Source:    "band",
			Timestamp: 110000,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	resp, err := suite.app.OracleKeeper.Price(sdk.WrapSDKContext(suite.ctx), &types.QueryGetPriceRequest{
		Asset: "BTC",
	})
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Price, prices[2])
	resp, err = suite.app.OracleKeeper.Price(sdk.WrapSDKContext(suite.ctx), &types.QueryGetPriceRequest{
		Asset:  "BTC",
		Source: "binance",
	})
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Price, prices[0])
	resp, err = suite.app.OracleKeeper.Price(sdk.WrapSDKContext(suite.ctx), &types.QueryGetPriceRequest{
		Asset:     "BTC",
		Source:    "binance",
		Timestamp: 100000,
	})
	suite.Require().NoError(err)
	suite.Require().Equal(resp.Price, prices[0])
	resp, err = suite.app.OracleKeeper.Price(sdk.WrapSDKContext(suite.ctx), &types.QueryGetPriceRequest{
		Asset:     "BTC",
		Source:    "binance",
		Timestamp: 11000,
	})
	suite.Require().Error(err)
}
