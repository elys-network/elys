package keeper_test

import (
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/clob/types"
	"time"
)

func (suite *KeeperTestSuite) TestTwapPrices() {

	market := types.PerpetualMarket{Id: 1, MaxTwapPricesTime: 20}
	suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)

	p1 := types.TwapPrice{
		MarketId:        1,
		Block:           1,
		Price:           math.LegacyNewDecWithPrec(101, 1),
		CumulativePrice: math.LegacyZeroDec(),
		Timestamp:       uint64(time.Now().Unix()),
	}
	p2 := types.TwapPrice{
		MarketId:        1,
		Block:           2,
		Price:           math.LegacyNewDecWithPrec(105, 1),
		CumulativePrice: math.LegacyZeroDec(),
		Timestamp:       p1.Timestamp + 5,
	}
	p3 := types.TwapPrice{
		MarketId:        1,
		Block:           3,
		Price:           math.LegacyNewDecWithPrec(111, 1),
		CumulativePrice: math.LegacyZeroDec(),
		Timestamp:       p2.Timestamp + 3,
	}
	p4 := types.TwapPrice{
		MarketId:        1,
		Block:           4,
		Price:           math.LegacyNewDecWithPrec(107, 1),
		CumulativePrice: math.LegacyZeroDec(),
		Timestamp:       p3.Timestamp + 17,
	}
	testCases := []struct {
		name   string
		result math.LegacyDec
		pre    func()
		post   func()
	}{
		{
			"first twap price is 0",
			math.LegacyZeroDec(),
			func() {
				suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p1) // cumulativePrice = 0
			},
			func() {
				lastMarketPrice := suite.app.ClobKeeper.GetLastMarketPrice(suite.ctx, 1)
				suite.Require().Equal(p1.Price, lastMarketPrice)
			},
		},
		{
			"first and 2nd price are same",
			math.LegacyNewDecWithPrec(101, 1), // (50.5 - 0)/5
			func() {
				suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p2) // cumulativePrice = 50.5
			},
			func() {
				lastMarketPrice := suite.app.ClobKeeper.GetLastMarketPrice(suite.ctx, 1)
				suite.Require().Equal(p2.Price, lastMarketPrice)
			},
		},
		{
			"third trade",
			math.LegacyNewDecWithPrec(1025, 2), // (82-0)/8
			func() {
				suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p3) // cumulativePrice = 50.5 + 10.5*3 = 82
			},
			func() {
				// Test GetAllTwapPrices
				allTwapPrices := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
				suite.Require().Len(allTwapPrices, 3)
				suite.Require().Equal(allTwapPrices[0], p1)
				p2.CumulativePrice = math.LegacyNewDecWithPrec(505, 1)
				suite.Require().Equal(allTwapPrices[1], p2)
				p3.CumulativePrice = math.LegacyNewDecWithPrec(82, 0)
				suite.Require().Equal(allTwapPrices[2], p3)

				lastMarketPrice := suite.app.ClobKeeper.GetLastMarketPrice(suite.ctx, 1)
				suite.Require().Equal(p3.Price, lastMarketPrice)
			},
		},
		{
			"should delete p1",
			math.LegacyNewDecWithPrec(1101, 2), // (270.7 - 50.5)/20
			func() {
				suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p4) // cumulativePrice = 82 + 11.1*17 = 270.7
			},
			func() {
				allTwapPrices := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
				suite.Require().Len(allTwapPrices, 3)
				suite.Require().Equal(allTwapPrices[0], p2)
				suite.Require().Equal(allTwapPrices[1], p3)
				p4.CumulativePrice = math.LegacyNewDecWithPrec(2707, 1)
				suite.Require().Equal(allTwapPrices[2], p4)

				lastMarketPrice := suite.app.ClobKeeper.GetLastMarketPrice(suite.ctx, 1)
				suite.Require().Equal(p4.Price, lastMarketPrice)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.pre()
			res := suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, 1)
			suite.Equal(tc.result.String(), res.String())
			tc.post()
		})
	}
}

func (suite *KeeperTestSuite) TestGetLowestSellPrice() {
	sell1 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1023, 2),
		BlockHeight: 1,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.NewInt(100),
		Filled:      math.ZeroInt(),
	}
	sell2 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1027, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.NewInt(100),
		Filled:      math.ZeroInt(),
	}
	sell3 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
		Price:       math.LegacyNewDecWithPrec(1029, 2),
		BlockHeight: 3,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.NewInt(100),
		Filled:      math.ZeroInt(),
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
		Amount:      math.NewInt(100),
		Filled:      math.ZeroInt(),
	}
	buy2 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1017, 2),
		BlockHeight: 2,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.NewInt(100),
		Filled:      math.ZeroInt(),
	}
	buy3 := types.PerpetualOrder{
		MarketId:    1,
		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
		Price:       math.LegacyNewDecWithPrec(1011, 2),
		BlockHeight: 3,
		Owner:       authtypes.NewModuleAddress("1").String(),
		Amount:      math.NewInt(100),
		Filled:      math.ZeroInt(),
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

//func (suite *KeeperTestSuite) TestGetMidPrice() {
//	sell1 := types.PerpetualOrder{
//		MarketId:    1,
//		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
//		Price:       math.LegacyNewDecWithPrec(1023, 2),
//		BlockHeight: 1,
//		Owner:       authtypes.NewModuleAddress("1").String(),
//		Amount:      math.NewInt(100),
//		Filled:      math.ZeroInt(),
//	}
//	sell2 := types.PerpetualOrder{
//		MarketId:    1,
//		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
//		Price:       math.LegacyNewDecWithPrec(1027, 2),
//		BlockHeight: 2,
//		Owner:       authtypes.NewModuleAddress("1").String(),
//		Amount:      math.NewInt(100),
//		Filled:      math.ZeroInt(),
//	}
//	sell3 := types.PerpetualOrder{
//		MarketId:    1,
//		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_SELL,
//		Price:       math.LegacyNewDecWithPrec(1029, 2),
//		BlockHeight: 3,
//		Owner:       authtypes.NewModuleAddress("1").String(),
//		Amount:      math.NewInt(100),
//		Filled:      math.ZeroInt(),
//	}
//
//	buy1 := types.PerpetualOrder{
//		MarketId:    1,
//		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
//		Price:       math.LegacyNewDecWithPrec(1013, 2),
//		BlockHeight: 1,
//		Owner:       authtypes.NewModuleAddress("1").String(),
//		Amount:      math.NewInt(100),
//		Filled:      math.ZeroInt(),
//	}
//	buy2 := types.PerpetualOrder{
//		MarketId:    1,
//		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
//		Price:       math.LegacyNewDecWithPrec(1017, 2),
//		BlockHeight: 2,
//		Owner:       authtypes.NewModuleAddress("1").String(),
//		Amount:      math.NewInt(100),
//		Filled:      math.ZeroInt(),
//	}
//	buy3 := types.PerpetualOrder{
//		MarketId:    1,
//		OrderType:   types.OrderType_ORDER_TYPE_LIMIT_BUY,
//		Price:       math.LegacyNewDecWithPrec(1011, 2),
//		BlockHeight: 3,
//		Owner:       authtypes.NewModuleAddress("1").String(),
//		Amount:      math.NewInt(100),
//		Filled:      math.ZeroInt(),
//	}
//
//	testCases := []struct {
//		name           string
//		result         math.LegacyDec
//		expectedErrMsg string
//		pre            func()
//	}{
//		{
//			"no mid price",
//			math.LegacyZeroDec(),
//			"one side of the orderbook is empty",
//			func() {
//				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy1)
//			},
//		},
//		{
//			"success 1",
//			buy2.Price,
//			"",
//			func() {
//				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sell1)
//			},
//		},
//		//{
//		//	"3rd price is set but price 10.17",
//		//	buy2.Price,
//		//	func() {
//		//		suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buy3)
//		//	},
//		//},
//	}
//
//	for _, tc := range testCases {
//		suite.Run(tc.name, func() {
//			tc.pre()
//			res, err := suite.app.ClobKeeper.GetMidPrice(suite.ctx, 1)
//			if e
//				suite.Equal(res.String(), tc.result.String())
//		})
//	}
//}
