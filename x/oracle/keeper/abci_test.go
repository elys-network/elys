package keeper_test

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// This test ensures old price data is automatically removed
func (suite *KeeperTestSuite) TestEndBlock() {
	now := time.Now()
	params := types.DefaultParams()
	suite.app.OracleKeeper.SetParams(suite.ctx, params)

	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(1),
			Source:    "binance",
			Timestamp: uint64(now.Unix()),
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(2),
			Source:    "band",
			Timestamp: uint64(now.Unix()) + params.PriceExpiryTime,
		},
		{
			Asset:     "BTC",
			Price:     sdk.NewDec(3),
			Source:    "band",
			Timestamp: uint64(now.Unix()) + params.PriceExpiryTime,
		},
	}
	for _, price := range prices {
		suite.app.OracleKeeper.SetPrice(suite.ctx, price)
	}
	suite.ctx = suite.ctx.WithBlockTime(now.Add(time.Second * time.Duration(params.PriceExpiryTime)))

	prices = suite.app.OracleKeeper.GetAllPrice(suite.ctx)
	suite.Require().Len(prices, 2)
}
