package keeper_test

import (
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestSetTwapPrices() {
	suite.ResetSuite()
	trade := types.Trade{
		Quantity: math.LegacyNewDec(-100),
	}
	err := suite.app.ClobKeeper.SetTwapPrices(suite.ctx, trade)
	suite.Require().Error(err)
	suite.Require().Equal("trade quantity cannot be negative or zero", err.Error())
}

func (suite *KeeperTestSuite) TestGetLowestSellPrice() {
	sell1 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1023, 2),
		BlockHeight: 1,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	sell2 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1027, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	sell3 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1029, 2),
		BlockHeight: 3,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}

	testCases := []struct {
		name   string
		result math.LegacyDec
		pre    func()
	}{
		{
			"first price 10.27",
			sell2.Price,
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell2)
			},
		},
		{
			"2nd price 10.23",
			sell1.Price,
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell1)
			},
		},
		{
			"3rd price is set but price 10.23",
			sell1.Price,
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell3)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.pre()
			res := suite.app.ClobKeeper.GetLowestSellPrice(suite.ctx, 1)
			suite.Equal(res.String(), tc.result.String())
		})
	}
}

func (suite *KeeperTestSuite) TestGetHighestBuyPrice() {
	buy1 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1013, 2),
		BlockHeight: 1,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	buy2 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1017, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	buy3 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1011, 2),
		BlockHeight: 3,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}

	testCases := []struct {
		name   string
		result math.LegacyDec
		pre    func()
	}{
		{
			"first price 10.13",
			buy1.Price,
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy1)
			},
		},
		{
			"2nd price 10.17",
			buy2.Price,
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy2)
			},
		},
		{
			"3rd price is set but price 10.17",
			buy2.Price,
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy3)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.pre()
			res := suite.app.ClobKeeper.GetHighestBuyPrice(suite.ctx, 1)
			suite.Equal(res.String(), tc.result.String())
		})
	}
}

func (suite *KeeperTestSuite) TestGetMidPrice() {
	sell1 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1023, 2),
		BlockHeight: 1,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	sell2 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1027, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	sell3 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1020, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}

	buy1 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1013, 2),
		BlockHeight: 1,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	buy2 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1011, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}
	buy3 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1020, 2),
		BlockHeight: 1,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.LegacyNewDec(100),
		Filled:      math.LegacyZeroDec(),
	}

	testCases := []struct {
		name           string
		result         math.LegacyDec
		expectedErrMsg string
		pre            func()
	}{
		{
			"no mid price",
			math.LegacyZeroDec(),
			"one side of the orderbook is empty",
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy1)
			},
		},
		{
			"success 1",
			math.LegacyNewDecWithPrec(1018, 2),
			"",
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell1)
			},
		},
		{
			"2nd sell order is set, but sell price is higher than previous, so no impact on mid price",
			math.LegacyNewDecWithPrec(1018, 2),
			"",
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell2)
			},
		},
		{
			"2nd buy order is set, but buy price is lower than previous, so no impact on mid price",
			math.LegacyNewDecWithPrec(1018, 2),
			"",
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy2)
			},
		},
		{
			"success 2, 3rd sell order with lower price",
			math.LegacyNewDecWithPrec(10165, 3),
			"",
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell3)
			},
		},
		{
			"success 2, 3rd buy order with higher price",
			math.LegacyNewDecWithPrec(102, 1),
			"",
			func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy3)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.pre()
			res, err := suite.app.ClobKeeper.GetMidPrice(suite.ctx, 1)
			if tc.expectedErrMsg != "" {
				suite.Equal(tc.expectedErrMsg, err.Error())
			} else {
				suite.Equal(res.String(), tc.result.String())
			}
		})
	}
}

func (suite *KeeperTestSuite) TestSetTwapPricesStruct() {
	suite.ResetSuite()
	suite.IncreaseHeight(1)
	p1 := types.TwapPrice{
		MarketId:          1,
		Block:             uint64(suite.ctx.BlockHeight()),
		AverageTradePrice: math.LegacyNewDecWithPrec(111, 1),
		TotalVolume:       math.LegacyNewDec(150),
		CumulativePrice:   math.LegacyNewDecWithPrec(1004, 1),
		Timestamp:         uint64(suite.ctx.BlockTime().Unix()),
	}
	suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, p1)
	all := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
	suite.Require().Equal(len(all), 1)
	suite.Require().Equal(all[0], p1)

	suite.IncreaseHeight(1)

	p2 := types.TwapPrice{
		MarketId:          1,
		Block:             uint64(suite.ctx.BlockHeight()),
		AverageTradePrice: math.LegacyNewDecWithPrec(120, 1),
		TotalVolume:       math.LegacyNewDec(120),
		CumulativePrice:   p1.CumulativePrice.Add(p1.AverageTradePrice.MulInt64(int64(suite.avgBlockTime))),
		Timestamp:         uint64(suite.ctx.BlockTime().Unix()),
	}
	suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, p2)
	all = suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
	suite.Require().Equal(len(all), 2)
	suite.Require().Equal(all[0], p1)
	suite.Require().Equal(all[1], p2)

	currentTwapPrice := suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, 1)
	suite.Require().Equal(currentTwapPrice, (p2.CumulativePrice.Sub(p1.CumulativePrice)).QuoInt64(int64(suite.avgBlockTime)))

	suite.IncreaseHeight(1)

	p3 := types.TwapPrice{
		MarketId:          1,
		Block:             uint64(suite.ctx.BlockHeight()),
		AverageTradePrice: math.LegacyNewDecWithPrec(125, 1),
		TotalVolume:       math.LegacyNewDec(160),
		CumulativePrice:   p2.CumulativePrice.Add(p2.AverageTradePrice.MulInt64(int64(suite.avgBlockTime))),
		Timestamp:         uint64(suite.ctx.BlockTime().Unix()),
	}
	suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, p3)
	all = suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
	suite.Require().Equal(len(all), 3)
	suite.Require().Equal(all[0], p1)
	suite.Require().Equal(all[1], p2)
	suite.Require().Equal(all[2], p3)

	currentTwapPrice = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, 1)
	suite.Require().Equal(currentTwapPrice, (p3.CumulativePrice.Sub(p1.CumulativePrice)).QuoInt64(int64(suite.avgBlockTime*2)))
}
