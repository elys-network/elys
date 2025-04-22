package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/clob/types"
)

func (suite *KeeperTestSuite) TestUpdateFundingRate() {
	suite.ResetSuite()

	markets := suite.CreateMarket(BaseDenom)
	market := markets[0]

	p1 := []types.Trade{
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(101, 1),
			Quantity: math.LegacyNewDec(100),
		},
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(103, 1),
			Quantity: math.LegacyNewDec(300),
		},
	}
	p2 := []types.Trade{
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(105, 1),
			Quantity: math.LegacyNewDec(200)},
		{
			MarketId: 1,
			Price:    math.LegacyNewDecWithPrec(110, 1),
			Quantity: math.LegacyNewDec(300)},
	}
	suite.IncreaseHeight(1)
	for _, p := range p1 {
		err := suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p) // cumulativePrice = 0
		suite.Require().NoError(err)
	}
	suite.IncreaseHeight(1)
	for _, p := range p2 {
		err := suite.app.ClobKeeper.SetTwapPrices(suite.ctx, p) // cumulativePrice = 10.25*5 = 51.25
		suite.Require().NoError(err)
	}
	twapPrice := math.LegacyMustNewDecFromStr("10.25")
	suite.Require().Equal(twapPrice, suite.app.ClobKeeper.GetCurrentTwapPrice(suite.ctx, market.Id))

	testCases := []struct {
		name           string
		expectedErrMsg string
		result         math.LegacyDec
		pre            func()
		post           func()
	}{
		{
			"oracle price not found",
			"asset price not found for denom (uatom)",
			math.LegacyDec{},
			func() {
				suite.app.OracleKeeper.RemovePrice(suite.ctx, "ATOM", "test", uint64(suite.ctx.BlockTime().Unix()-15))
			},
			func() {
			},
		},
		{
			"funding rate increases within limit",
			"",
			math.LegacyMustNewDecFromStr("0.005025125628140704"),
			func() {
				oraclePrice := twapPrice.Sub(twapPrice.Mul(market.MaxAbsFundingRateChange.Quo(math.LegacyNewDec(2)))) //twap price(1 - MaxAbsFundingRateChange/2) = 10.19875
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			"funding rate decreases within limit",
			"",
			math.LegacyMustNewDecFromStr("-0.002493765586034913"),
			func() {
				oraclePrice := twapPrice.Add(twapPrice.Mul(market.MaxAbsFundingRateChange.Quo(math.LegacyNewDec(4)))) //twap price(1 + MaxAbsFundingRateChange/4) = 10.275625
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			"funding rate increases exceeding limit",
			"",
			math.LegacyMustNewDecFromStr("-0.002493765586034913").Add(market.MaxAbsFundingRateChange),
			func() {
				oraclePrice := twapPrice.Sub(twapPrice.Mul(market.MaxAbsFundingRateChange.MulInt64(2))) //twap price(1 - 2*MaxAbsFundingRateChange)
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			"funding rate decreases exceeding limit",
			"",
			math.LegacyMustNewDecFromStr("-0.002493765586034913"),
			func() {
				oraclePrice := twapPrice.Add(twapPrice.Mul(market.MaxAbsFundingRateChange.MulInt64(2))) //twap price(1 - 2*MaxAbsFundingRateChange)
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{oraclePrice, math.LegacyNewDec(1)})
			},
			func() {
			},
		},
		{
			name: "Rate hits positive absolute cap",
			// Expect: lastRate(0.0095) + change(clamped to +0.001) = 0.0105 -> clamped to MaxAbsRate(0.01)
			result: market.MaxAbsFundingRate, // Assuming market.MaxAbsFundingRate = 0.01
			pre: func() {
				// Set last rate just below the absolute cap
				lastRate := market.MaxAbsFundingRate.Sub(market.MaxAbsFundingRateChange.QuoInt64(100)) // e.g., 0.01 - 0.001/2 = 0.0095
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{
					MarketId: market.Id,
					Block:    uint64(suite.ctx.BlockHeight()),
					Rate:     lastRate,
				})
				// Set index price significantly lower than TWAP to force large positive change
				indexPrice := math.LegacyMustNewDecFromStr("5")
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{indexPrice, math.LegacyNewDec(1)}) // Assuming BaseDenom is "ATOM" from previous example
			},
			post: func() {},
		},
		{
			name: "Rate hits negative absolute cap",
			// Expect: lastRate(-0.0095) + change(clamped to -0.001) = -0.0105 -> clamped to -MaxAbsRate(-0.01)
			result: market.MaxAbsFundingRate.Neg(), // Assuming market.MaxAbsFundingRate = 0.01
			pre: func() {
				// Set last rate just above the negative absolute cap
				lastRate := market.MaxAbsFundingRate.Sub(market.MaxAbsFundingRateChange.QuoInt64(100)).Neg() // e.g., -0.0095
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{
					MarketId: market.Id,
					Block:    uint64(suite.ctx.BlockHeight()),
					Rate:     lastRate,
				})
				// Set index price significantly higher than TWAP to force large negative change
				indexPrice := math.LegacyNewDec(20)
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{indexPrice, math.LegacyNewDec(1)})

			},
			post: func() {},
		},
		{
			name: "Rate calculation with zero TWAP price",
			// Expect: twap=0, index=10. premium=-10. rateCal=-1. lastRate=0.025 change=-1.005. Clamp change to -0.01. newRate=0.025-0.01=0.015
			result: math.LegacyMustNewDecFromStr("0.015"), // Adjust based on MaxAbsFundingRateChange and chosen lastRate
			pre: func() {
				// Clear previous TWAP data (needs helper)
				suite.ResetSuite()
				markets = suite.CreateMarket(BaseDenom)
				market = markets[0]
				// Set a last funding rate
				lastRate := market.MaxAbsFundingRate.QuoInt64(2) // e.g., 0.005
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{
					MarketId: market.Id,
					Block:    uint64(suite.ctx.BlockHeight()),
					Rate:     lastRate,
				})
				// Set a valid index price
				indexPrice := math.LegacyNewDec(10)
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{indexPrice, math.LegacyNewDec(1)})
			},
			post: func() {
				// Optional: Assert GetCurrentTwapPrice was indeed zero if ClearTwapData isn't guaranteed
				// suite.Require().True(suite.keeper.GetCurrentTwapPrice(suite.ctx, market.Id).IsZero())
			},
		},
		{
			name:   "Rate change with zero premium (change clamp hit)",
			result: math.LegacyMustNewDecFromStr("0.01"), // Adjust based on MaxAbsFundingRateChange
			pre: func() {
				// Set last rate such that abs(0 - lastRate) > MaxAbsFundingRateChange
				lastRate := market.MaxAbsFundingRateChange.MulInt64(2) // e.g., 0.002 if MaxChange is 0.001
				// Ensure lastRate itself is within absolute limits if necessary for setup
				if lastRate.Abs().GT(market.MaxAbsFundingRate) {
					lastRate = market.MaxAbsFundingRate
				}
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{
					MarketId: market.Id,
					Block:    uint64(suite.ctx.BlockHeight()),
					Rate:     lastRate,
				})
				// Set index price == base twap price
				indexPrice := math.LegacyMustNewDecFromStr("10.25") // Matches TWAP from test setup
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{indexPrice, math.LegacyNewDec(1)})
			},
			post: func() {},
		},
		{
			name: "Rate change with zero premium (no change clamp hit)",
			// Expect: twap=10.25, index=10.25. premium=0. rateCal=0. lastRate=0.0005 (example < MaxChange). change=0-0.0005=-0.0005. Clamp change? No. newRate=0.0005-0.0005=0
			result: math.LegacyMustNewDecFromStr("-0.005"),
			pre: func() {
				suite.ResetSuite()
				markets = suite.CreateMarket(BaseDenom)
				market = markets[0]
				// Set last rate such that abs(0 - lastRate) <= MaxAbsFundingRateChange
				lastRate := market.MaxAbsFundingRateChange.QuoInt64(2) // e.g., 0.0005 if MaxChange is 0.001
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{
					MarketId: market.Id,
					Block:    uint64(suite.ctx.BlockHeight()),
					Rate:     lastRate,
				})
				// Set index price == base twap price
				indexPrice := math.LegacyMustNewDecFromStr("10.25")
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{indexPrice, math.LegacyNewDec(1)})
			},
			post: func() {},
		},
		{
			name:           "Error on zero index price",
			expectedErrMsg: "asset price",    // Expect error from GetAssetPrice containing this
			result:         math.LegacyDec{}, // Result is irrelevant on error
			pre: func() {
				// Set oracle price to zero
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{math.LegacyZeroDec(), math.LegacyNewDec(1)})
			},
			post: func() {},
		},
		{
			name:           "Error on negative index price",
			expectedErrMsg: "asset price",    // Expect error from GetAssetPrice containing this
			result:         math.LegacyDec{}, // Result is irrelevant on error
			pre: func() {
				// Set oracle price to negative
				suite.SetPrice([]string{"ATOM", "USDC"}, []math.LegacyDec{math.LegacyNewDec(-10), math.LegacyNewDec(1)})
			},
			post: func() {},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			suite.IncreaseHeight(1)
			tc.pre()
			err := suite.app.ClobKeeper.UpdateFundingRate(suite.ctx, market)
			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
				fundingRate := suite.app.ClobKeeper.GetFundingRate(suite.ctx, market.Id)
				suite.Require().Equal(tc.result, fundingRate.Rate)
				suite.Require().Equal(uint64(suite.ctx.BlockHeight()), fundingRate.Block)
			}
			tc.post()
		})
	}

}

func (suite *KeeperTestSuite) TestFundingRate() {
	suite.ResetSuite()
	markets := suite.CreateMarket(BaseDenom, "uosmo")
	suite.SetPrice([]string{"OSMO"}, []math.LegacyDec{math.LegacyNewDec(2)})
	err := suite.app.ClobKeeper.UpdateFundingRate(suite.ctx, markets[0])
	suite.Require().NoError(err)
	err = suite.app.ClobKeeper.UpdateFundingRate(suite.ctx, markets[1])
	suite.Require().NoError(err)

	all := suite.app.ClobKeeper.GetAllFundingRate(suite.ctx)
	suite.Require().Len(all, 2)

	suite.Require().Equal(math.LegacyMustNewDecFromStr("-0.01"), all[0].Rate)
	suite.Require().Equal(uint64(suite.ctx.BlockHeight()), all[0].Block)
	suite.Require().Equal(uint64(1), all[0].MarketId)

	suite.Require().Equal(math.LegacyMustNewDecFromStr("-0.01"), all[1].Rate)
	suite.Require().Equal(uint64(suite.ctx.BlockHeight()), all[1].Block)
	suite.Require().Equal(uint64(2), all[1].MarketId)
}
