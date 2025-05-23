package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"

	"github.com/elys-network/elys/v5/x/oracle/types"
)

// This test ensures old price data is automatically removed
func (suite *KeeperTestSuite) TestEndBlock() {
	now := time.Now()
	params := types.DefaultParams()
	suite.app.OracleKeeper.SetParams(suite.ctx, params)

	prices := []types.Price{
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(1),
			Source:    "elys",
			Timestamp: uint64(now.Unix()),
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(2),
			Source:    "band",
			Timestamp: uint64(now.Unix()) + params.PriceExpiryTime,
		},
		{
			Asset:     "BTC",
			Price:     sdkmath.LegacyNewDec(3),
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
