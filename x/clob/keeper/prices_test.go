package keeper_test

import (
	"cosmossdk.io/math"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/clob/types"
	"time"
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

func (suite *KeeperTestSuite) TestSetTwapFunctions() {
	// --- Setup ---
	var market types.PerpetualMarket
	var buyerAcc, sellerAcc types.SubAccount       // Dummy accounts for trade obj
	var tradeTimeStep = AvgBlockTime * time.Second // Simulate 5 seconds per block height increase
	var twapWindowSeconds uint64 = 15              // Match MaxTwapPricesTime in original test

	setupTest := func() {
		_, buyerAcc, sellerAcc, _ = suite.SetupExchangeTest() // Re-setup accounts/market state
		market = types.PerpetualMarket{Id: MarketId, TwapPricesWindow: twapWindowSeconds}
		suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market) // Ensure market with window is set
		// Reset context time maybe? Assuming SetupExchangeTest provides fresh ctx
		suite.ctx = suite.ctx.WithBlockTime(time.Unix(1700000000, 0)) // Set a known start time
		suite.ctx = suite.ctx.WithBlockHeight(100)                    // Set known start height
	}

	p1 := []types.Trade{ // Trades for Block H+1 (Time T0+5s)
		types.NewTrade(MarketId, math.LegacyNewDec(100), math.LegacyMustNewDecFromStr("10.1"), buyerAcc, sellerAcc, false),
		types.NewTrade(MarketId, math.LegacyNewDec(300), math.LegacyMustNewDecFromStr("10.3"), buyerAcc, sellerAcc, false),
	} // AvgPrice = (1010 + 3090) / 400 = 10.25
	p1AvgPrice := math.LegacyMustNewDecFromStr("10.25")

	p2 := []types.Trade{ // Trades for Block H+2 (Time T0+10s)
		types.NewTrade(MarketId, math.LegacyNewDec(200), math.LegacyMustNewDecFromStr("10.5"), buyerAcc, sellerAcc, false),
		types.NewTrade(MarketId, math.LegacyNewDec(300), math.LegacyMustNewDecFromStr("11.0"), buyerAcc, sellerAcc, false),
	} // AvgPrice = (2100 + 3300) / 500 = 10.8
	p2AvgPrice := math.LegacyMustNewDecFromStr("10.8")

	//p3 := []types.Trade{ // Trades for Block H+3 (Time T0+15s)
	//	types.NewTrade(MarketId, math.LegacyNewDec(100), math.LegacyMustNewDecFromStr("11.1"), buyerAcc, sellerAcc, false),
	//	types.NewTrade(MarketId, math.LegacyNewDec(300), math.LegacyMustNewDecFromStr("10.9"), buyerAcc, sellerAcc, false),
	//} // AvgPrice = (1110 + 3270) / 400 = 10.95
	p3AvgPrice := math.LegacyMustNewDecFromStr("10.95")

	p4 := []types.Trade{ // Trades for Block H+5 (Time T0+25s) - NOTE Height increases by 2 here
		types.NewTrade(MarketId, math.LegacyNewDec(300), math.LegacyMustNewDecFromStr("11.5"), buyerAcc, sellerAcc, false),
		types.NewTrade(MarketId, math.LegacyNewDec(200), math.LegacyMustNewDecFromStr("11.6"), buyerAcc, sellerAcc, false),
	} // AvgPrice = (3450 + 2320) / 500 = 11.54
	//p4AvgPrice := math.LegacyMustNewDecFromStr("11.54")
	p100 := math.LegacyMustNewDecFromStr("100")
	// === Test GetCurrentTwapPrice Logic ===
	p5 := []types.Trade{ // Trades for Block H+5 (Time T0+25s) - NOTE Height increases by 2 here
		types.NewTrade(2, math.LegacyNewDec(300), math.LegacyMustNewDecFromStr("11.5"), buyerAcc, sellerAcc, false),
	}
	suite.Run("GetCurrentTwapPrice Scenarios", func() {
		setupTest()
		market, _, _, _ = suite.SetupExchangeTest() // Re-get market with correct window

		// 1. No Records
		twap := suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		suite.Require().True(twap.IsZero(), "TWAP should be zero with no records")

		// 2. One Record
		suite.IncreaseHeight(1) // H+1, T0+5s
		currentTime := uint64(suite.ctx.BlockTime().Unix())
		currentBlock := uint64(suite.ctx.BlockHeight())
		suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
			MarketId: MarketId, Block: currentBlock, AverageTradePrice: p1AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: math.LegacyZeroDec(), Timestamp: currentTime,
		})
		twap = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		// Should return the average price of the single record
		suite.Require().True(twap.Equal(p1AvgPrice), "TWAP should equal single record avg price. Exp %s, Got %s", p1AvgPrice, twap)

		// 3. Two Records
		suite.IncreaseHeight(1) // H+2, T0+10s
		currentTime2 := uint64(suite.ctx.BlockTime().Unix())
		currentBlock2 := uint64(suite.ctx.BlockHeight())
		cumPrice2 := math.LegacyZeroDec().Add(p1AvgPrice.MulInt64(5)) // 0 + 10.25*5 = 51.25
		suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
			MarketId: MarketId, Block: currentBlock2, AverageTradePrice: p2AvgPrice, TotalVolume: math.LegacyNewDec(500), CumulativePrice: cumPrice2, Timestamp: currentTime2,
		})
		twap = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		// TWAP = (CumPrice2 - CumPrice1) / (Time2 - Time1) = (51.25 - 0) / 5 = 10.25
		expectedTwap := math.LegacyMustNewDecFromStr("10.25")
		suite.Require().True(expectedTwap.Equal(twap), "TWAP mismatch with two records. Exp %s, Got %s", expectedTwap, twap)

		// 4. Three Records
		suite.IncreaseHeight(1) // H+3, T0+15s
		currentTime3 := uint64(suite.ctx.BlockTime().Unix())
		currentBlock3 := uint64(suite.ctx.BlockHeight())
		cumPrice3 := cumPrice2.Add(p2AvgPrice.MulInt64(5)) // 51.25 + 10.8*5 = 51.25 + 54 = 105.25
		suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
			MarketId: MarketId, Block: currentBlock3, AverageTradePrice: p3AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: cumPrice3, Timestamp: currentTime3,
		})
		twap = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		// TWAP = (CumPrice3 - CumPrice1) / (Time3 - Time1) = (105.25 - 0) / 10 = 10.525
		expectedTwap = math.LegacyMustNewDecFromStr("10.525")
		suite.Require().True(expectedTwap.Equal(twap), "TWAP mismatch with three records. Exp %s, Got %s", expectedTwap, twap)

		// 5. Test Timestamp Panic Condition
		suite.Run("Panic on inverted timestamps", func() {
			// Manually set a bad record where last time < first time
			badFirstTime := uint64(suite.ctx.BlockTime().Unix() + 100) // Future time
			badLastTime := uint64(suite.ctx.BlockTime().Unix())        // Current time
			suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: 1, AverageTradePrice: math.LegacyOneDec(), CumulativePrice: math.LegacyZeroDec(), Timestamp: badFirstTime})
			suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: 2, AverageTradePrice: math.LegacyOneDec(), CumulativePrice: math.LegacyOneDec(), Timestamp: badLastTime})

			call := func() {
				_ = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
			}
			suite.Require().PanicsWithValue("twap price timestamp delta incorrect, time delta < 0", call)
		})
	})

	// === Test SetTwapPrices Logic and Errors ===
	suite.Run("SetTwapPrices Scenarios", func() {
		testCases := []struct {
			name        string
			trade       types.Trade
			setupFunc   func() // Setup prior state if needed
			checkFunc   func() // Checks after successful call
			expectedErr string // Error substring, "" for no error
		}{
			{
				name:        "Error: Zero Quantity Trade",
				trade:       types.Trade{buyerAcc, sellerAcc, MarketId, p100, math.LegacyZeroDec(), false}, // Zero Qty
				setupFunc:   func() { setupTest() },                                                        // Basic setup
				expectedErr: "trade quantity cannot be negative or zero",
			},
			{
				name:        "Error: Negative Quantity Trade",
				trade:       types.Trade{buyerAcc, sellerAcc, MarketId, p100, math.LegacyNewDec(-10), false}, // Negative Qty
				setupFunc:   func() { setupTest() },
				expectedErr: "trade quantity cannot be negative or zero",
			},
			{
				name:      "First trade in block (Genesis)",
				trade:     p1[0], // 10.1, 100
				setupFunc: func() { setupTest() },
				checkFunc: func() {
					twap := suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
					suite.Require().True(twap.Equal(p1[0].Price), "TWAP should be first trade price") // Single record returns its avg price
					allTwaps := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)                      // Assuming helper exists
					suite.Require().Len(allTwaps, 1)
					suite.Require().True(allTwaps[0].AverageTradePrice.Equal(p1[0].Price))
					suite.Require().True(allTwaps[0].CumulativePrice.IsZero()) // First record has zero cumulative price start
					suite.Require().Equal(uint64(suite.ctx.BlockHeight()), allTwaps[0].Block)
				},
				expectedErr: "",
			},
			{
				name:  "Second trade in same block",
				trade: p1[1], // 10.3, 300
				setupFunc: func() {
					setupTest()
					// Manually set state as if p1[0] already happened in this block
					block := uint64(suite.ctx.BlockHeight())
					timeStamp := uint64(suite.ctx.BlockTime().Unix())
					initialTwap := types.TwapPrice{
						MarketId: MarketId, Block: block, Timestamp: timeStamp,
						AverageTradePrice: p1[0].Price, TotalVolume: p1[0].Quantity, CumulativePrice: math.LegacyZeroDec(),
					}
					suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, initialTwap)
				},
				checkFunc: func() {
					twap := suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
					suite.Require().True(twap.Equal(p1AvgPrice), "TWAP should be avg price of the single record") // Still only one record
					allTwaps := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
					suite.Require().Len(allTwaps, 1)
					suite.Require().True(allTwaps[0].AverageTradePrice.Equal(p1AvgPrice))                   // Check block avg price updated
					suite.Require().True(allTwaps[0].TotalVolume.Equal(p1[0].Quantity.Add(p1[1].Quantity))) // Check block volume updated
					suite.Require().True(allTwaps[0].CumulativePrice.IsZero())                              // Cumulative price not updated intra-block
				},
				expectedErr: "",
			},
			{
				name:  "First trade in new block (calculates cumulative)",
				trade: p2[0], // 10.5, 200
				setupFunc: func() {
					setupTest()
					// Set state for previous block (Block H-1, Time T0-5s)
					prevBlock := uint64(suite.ctx.BlockHeight() - 1)
					prevTime := uint64(suite.ctx.BlockTime().Unix() - int64(tradeTimeStep.Seconds()))
					prevTwap := types.TwapPrice{
						MarketId: MarketId, Block: prevBlock, Timestamp: prevTime,
						AverageTradePrice: p1AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: math.LegacyZeroDec(),
					}
					suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, prevTwap)
				},
				checkFunc: func() {
					allTwaps := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
					suite.Require().Len(allTwaps, 2)  // Should have prev and current block records
					currentBlockRecord := allTwaps[1] // Assuming order
					if allTwaps[0].Block > allTwaps[1].Block {
						currentBlockRecord = allTwaps[0]
					} // Find current block record
					suite.Require().Equal(uint64(suite.ctx.BlockHeight()), currentBlockRecord.Block)
					suite.Require().True(currentBlockRecord.AverageTradePrice.Equal(p2[0].Price)) // Avg price is just this trade's price
					suite.Require().True(currentBlockRecord.TotalVolume.Equal(p2[0].Quantity))    // Volume is just this trade's volume
					// Check cumulative: prevCum + prevAvg * timeDelta = 0 + 10.25 * 5 = 51.25
					expectedCum := math.LegacyMustNewDecFromStr("51.25")
					suite.Require().True(expectedCum.Equal(currentBlockRecord.CumulativePrice), "Cumulative price mismatch. Exp %s, Got %s", expectedCum, currentBlockRecord.CumulativePrice)
				},
				expectedErr: "",
			},
			{
				name:  "Pruning deletes old record",
				trade: p4[0], // 11.5, 300 - Trade in block H+5 (T0+25s)
				setupFunc: func() {
					setupTest()
					// Set records for H+1, H+2, H+3 (Times T0+5, T0+10, T0+15)
					t0 := suite.ctx.BlockTime()
					h0 := suite.ctx.BlockHeight()
					// H+1 / T0+5s
					suite.ctx = suite.ctx.WithBlockHeight(h0 + 1).WithBlockTime(t0.Add(1 * tradeTimeStep))
					cum1 := math.LegacyZeroDec()
					suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: uint64(h0 + 1), Timestamp: uint64(suite.ctx.BlockTime().Unix()), AverageTradePrice: p1AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: cum1})
					// H+2 / T0+10s
					suite.ctx = suite.ctx.WithBlockHeight(h0 + 2).WithBlockTime(t0.Add(2 * tradeTimeStep))
					cum2 := cum1.Add(p1AvgPrice.MulInt64(5)) // 51.25
					suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: uint64(h0 + 2), Timestamp: uint64(suite.ctx.BlockTime().Unix()), AverageTradePrice: p2AvgPrice, TotalVolume: math.LegacyNewDec(500), CumulativePrice: cum2})
					// H+3 / T0+15s
					suite.ctx = suite.ctx.WithBlockHeight(h0 + 3).WithBlockTime(t0.Add(3 * tradeTimeStep))
					cum3 := cum2.Add(p2AvgPrice.MulInt64(5)) // 105.25
					suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: uint64(h0 + 3), Timestamp: uint64(suite.ctx.BlockTime().Unix()), AverageTradePrice: p3AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: cum3})
					// Set context for the actual trade call (H+5 / T0+25s)
					suite.ctx = suite.ctx.WithBlockHeight(h0 + 5).WithBlockTime(t0.Add(5 * tradeTimeStep))
				},
				checkFunc: func() {
					allTwaps := suite.app.ClobKeeper.GetAllTwapPrices(suite.ctx)
					suite.Require().Len(allTwaps, 3, "Expected 3 records after pruning, got %d", len(allTwaps)) // H+2, H+3, H+5 should remain
					// Check timestamps if possible, or just check blocks
					blocksFound := make(map[int64]bool)
					for _, p := range allTwaps {
						blocksFound[int64(p.Block)] = true
					}
					suite.Require().False(blocksFound[suite.ctx.BlockHeight()-4], "Record from H+1 should have been pruned") // H+1 is 4 blocks before H+5
					suite.Require().True(blocksFound[suite.ctx.BlockHeight()-3], "Record from H+2 should exist")
					suite.Require().True(blocksFound[suite.ctx.BlockHeight()-2], "Record from H+3 should exist")
					suite.Require().True(blocksFound[int64(suite.ctx.BlockHeight())], "Record from current block H+5 should exist") // H+5 is current block
				},
				expectedErr: "",
			},
			{
				name:  "Error: Invalid Timestamp Delta (Current <= Last)",
				trade: p2[0],
				setupFunc: func() {
					setupTest()
					// Set a previous record with a timestamp *in the future* relative to context
					prevBlock := uint64(suite.ctx.BlockHeight() - 1)
					futureTime := uint64(suite.ctx.BlockTime().Unix() + 100) // Time > current ctx time
					prevTwap := types.TwapPrice{
						MarketId: MarketId, Block: prevBlock, Timestamp: futureTime, // Future timestamp
						AverageTradePrice: p1AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: math.LegacyZeroDec(),
					}
					suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, prevTwap)
				},
				checkFunc:   func() {}, // No checks needed after error
				expectedErr: "twap price timestamp delta incorrect",
			},
			{
				name:  "Error: Market Not Found During Pruning",
				trade: p5[0], // Trigger else block to reach pruning
				setupFunc: func() {
					setupTest()
				},
				checkFunc:   func() {},
				expectedErr: types.ErrPerpetualMarketNotFound.Error(), // Expect market not found from GetPerpetualMarket inside SetTwapPrices
			},
		}

		for _, tc := range testCases {
			suite.Run(tc.name, func() {
				tc.setupFunc() // Run setup specific to the test case

				// Execute SetTwapPrices
				err := suite.app.ClobKeeper.SetTwapPrices(suite.ctx, tc.trade)

				// Assert Error
				if tc.expectedErr != "" {
					suite.Require().Error(err, "Expected an error but got nil")
					suite.Require().Contains(err.Error(), tc.expectedErr, "Error message mismatch")
				} else {
					suite.Require().NoError(err, "Expected no error but got: %v", err)
					// Run post-success checks if defined
					if tc.checkFunc != nil {
						tc.checkFunc()
					}
				}
			})
		}
	}) // End SetTwapPrices Scenarios

	// Add the GetCurrentTwapPrice Scenarios again to ensure they use the LATEST version of GetCurrentTwapPrice
	suite.Run("GetCurrentTwapPrice Scenarios (Post-Fix)", func() {
		// ... (copy the GetCurrentTwapPrice tests from previous response, ensuring they use the fixed version) ...
		setupTest()
		market, _ = suite.app.ClobKeeper.GetPerpetualMarket(suite.ctx, MarketId) // Re-get market with correct window

		// 1. No Records
		twap := suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		suite.Require().True(twap.IsZero(), "TWAP should be zero with no records")

		// 2. One Record
		suite.IncreaseHeight(1) // H+1, T0+5s
		currentTime := uint64(suite.ctx.BlockTime().Unix())
		currentBlock := uint64(suite.ctx.BlockHeight())
		suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
			MarketId: MarketId, Block: currentBlock, AverageTradePrice: p1AvgPrice, TotalVolume: math.LegacyNewDec(400), CumulativePrice: math.LegacyZeroDec(), Timestamp: currentTime,
		})
		twap = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		// Should return the average price of the single record
		suite.Require().True(twap.Equal(p1AvgPrice), "TWAP should equal single record avg price. Exp %s, Got %s", p1AvgPrice, twap)

		// 3. Two Records
		suite.IncreaseHeight(1) // H+2, T0+10s
		currentTime2 := uint64(suite.ctx.BlockTime().Unix())
		currentBlock2 := uint64(suite.ctx.BlockHeight())
		cumPrice2 := math.LegacyZeroDec().Add(p1AvgPrice.MulInt64(int64(tradeTimeStep.Seconds()))) // Use duration
		suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{
			MarketId: MarketId, Block: currentBlock2, AverageTradePrice: p2AvgPrice, TotalVolume: math.LegacyNewDec(500), CumulativePrice: cumPrice2, Timestamp: currentTime2,
		})
		twap = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
		// TWAP = (CumPrice2 - CumPrice1) / (Time2 - Time1) = (51.25 - 0) / 5 = 10.25
		expectedTwap := math.LegacyMustNewDecFromStr("10.25")
		suite.Require().True(expectedTwap.Equal(twap), "TWAP mismatch with two records. Exp %s, Got %s", expectedTwap, twap)

		// 4. Test Timestamp Panic Condition
		suite.Run("Panic on inverted timestamps", func() {
			// Need fresh setup to avoid state pollution
			setupTest()
			market, _, _, _ = suite.SetupExchangeTest()
			now := suite.ctx.BlockTime()
			// Manually set a bad record where last time < first time
			badFirstTime := uint64(now.Unix() + 100) // Future time
			badLastTime := uint64(now.Unix())        // Current time
			suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: 1, AverageTradePrice: math.LegacyOneDec(), CumulativePrice: math.LegacyZeroDec(), Timestamp: badFirstTime})
			suite.app.ClobKeeper.SetTwapPricesStruct(suite.ctx, types.TwapPrice{MarketId: MarketId, Block: 2, AverageTradePrice: math.LegacyOneDec(), CumulativePrice: math.LegacyOneDec(), Timestamp: badLastTime})

			call := func() {
				_ = suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, MarketId)
			}
			// Assuming the panic remains, test for it
			suite.Require().PanicsWithValue("twap price timestamp delta incorrect, time delta < 0", call)
		})

	}) // End GetCurrentTwapPrice Scenarios

}
