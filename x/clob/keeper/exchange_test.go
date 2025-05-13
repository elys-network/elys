package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/clob/types"
	"time"
)

func (suite *KeeperTestSuite) TestExchange_Comprehensive() {
	// Define common values 	// 10%
	denomFactor := math.NewInt(1000_000)

	calcMargin := func(qty, price math.LegacyDec) math.Int { // Local helper
		if price.IsNil() || !price.IsPositive() || qty.IsZero() {
			return math.ZeroInt()
		}
		return qty.Abs().Mul(price).Mul(IMR).MulInt(denomFactor).RoundInt()
	}
	p100 := math.LegacyNewDec(100)
	p105 := math.LegacyNewDec(105)
	p95 := math.LegacyNewDec(95)
	qty5 := math.LegacyNewDec(5)
	qty10 := math.LegacyNewDec(10)
	qty15 := math.LegacyNewDec(15)

	type ExpectedState struct {
		Err                  string
		BuyerExists          bool
		SellerExists         bool
		BuyerFinalQty        math.LegacyDec
		BuyerFinalEP         math.LegacyDec
		BuyerFinalMargin     math.Int
		SellerFinalQty       math.LegacyDec
		SellerFinalEP        math.LegacyDec
		SellerFinalMargin    math.Int
		BuyerBalanceChange   math.Int
		SellerBalanceChange  math.Int
		MarketBalanceChange  math.Int
		FinalMarketOI        math.LegacyDec
		TwapBlockDataSet     bool
		BuyerFundingRateSet  bool
		SellerFundingRateSet bool
	}

	testCases := []struct {
		name     string
		setup    func() (buyerAcc types.SubAccount, sellerAcc types.SubAccount)
		trade    func(buyer types.SubAccount, seller types.SubAccount) types.Trade
		expected ExpectedState
	}{
		// === GROUP 1: Both Start Empty ===
		{
			name: "B(0), S(0): Buyer Opens Long, Seller Opens Short",
			setup: func() (types.SubAccount, types.SubAccount) {
				_, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty10, p100, b, s, false) }, // B buys 10, S sells 10 @ 100
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: true,
				BuyerFinalQty: qty10, BuyerFinalEP: p100, BuyerFinalMargin: math.NewInt(100).Mul(denomFactor), // Q=10, EP=100, M=100
				SellerFinalQty: qty10.Neg(), SellerFinalEP: p100, SellerFinalMargin: math.NewInt(100).Mul(denomFactor), // Q=-10, EP=100, M=100
				BuyerBalanceChange: math.NewInt(-100).Mul(denomFactor), SellerBalanceChange: math.NewInt(-100).Mul(denomFactor), MarketBalanceChange: math.NewInt(200).Mul(denomFactor),
				FinalMarketOI: qty10, TwapBlockDataSet: true, BuyerFundingRateSet: true, SellerFundingRateSet: true,
			},
		},

		// === GROUP 2: One Starts Empty, One Has Position ===
		// B(0) S(+) -> B Open L, S Dec L
		{
			name: "B(0), S(+): Buyer Opens Long, Seller Decreases Long",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				sellerPerp := types.NewPerpetual(0, MarketId, sellerAcc.Owner, qty15, p100, calcMargin(qty15, p100), math.LegacyZeroDec()) // Seller +15 @ 100 (M=150)
				suite.SetPerpetualStateWithEntryFR(sellerPerp)
				market.TotalOpen = qty15
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty10, p105, b, s, false) }, // B buys 10, S sells 10 @ 105
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: true,
				BuyerFinalQty: qty10, BuyerFinalEP: p105, BuyerFinalMargin: calcMargin(qty10, p105), // Q=10, EP=105, M=105
				SellerFinalQty: qty5, SellerFinalEP: p100, SellerFinalMargin: calcMargin(qty5, p100), // Q=5, EP=100, M=50
				BuyerBalanceChange: math.NewInt(-105).Mul(denomFactor), // -Margin
				// Seller: PNL=+10*(105-100)=+50. Margin Refund = 150-50 = +100. Total=+150
				SellerBalanceChange: math.NewInt(150).Mul(denomFactor),
				MarketBalanceChange: math.NewInt(-45).Mul(denomFactor), // +B_M(105) - S_PNL(50) - S_Refund(100) = -45
				FinalMarketOI:       qty15,                             // Unchanged (1 Open, 1 Dec) = 15
				TwapBlockDataSet:    true, BuyerFundingRateSet: true, SellerFundingRateSet: true,
			},
		},
		// B(0) S(+) -> B Open L, S Close L
		{
			name: "B(0), S(+): Buyer Opens Long, Seller Closes Long",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				sellerPerp := types.NewPerpetual(0, MarketId, sellerAcc.Owner, qty10, p100, calcMargin(qty10, p100), math.LegacyZeroDec()) // Seller +10 @ 100 (M=100)
				suite.SetPerpetualStateWithEntryFR(sellerPerp)
				market.TotalOpen = qty10
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty10, p105, b, s, false) }, // B buys 10, S sells 10 @ 105
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: false, // Seller closed
				BuyerFinalQty: qty10, BuyerFinalEP: p105, BuyerFinalMargin: calcMargin(qty10, p105), // Q=10, EP=105, M=105
				// SellerFinalPerpetual: N/A
				BuyerBalanceChange: math.NewInt(-105).Mul(denomFactor), // -Margin
				// Seller PNL = +10*(105-100)=+50. Refund=100. Total=+150
				SellerBalanceChange: math.NewInt(150).Mul(denomFactor),
				MarketBalanceChange: math.NewInt(-45).Mul(denomFactor),                            // +B_M(105) - S_PNL(50) - S_Refund(100) = -45
				FinalMarketOI:       qty10,                                                        // Unchanged (1 Open, 1 Close) = 10
				TwapBlockDataSet:    true, BuyerFundingRateSet: true, SellerFundingRateSet: false, // Seller deleted
			},
		},
		// B(0) S(+) -> B Open L, S Flip L->S
		{
			name: "B(0), S(+): Buyer Opens Long, Seller Flips Long->Short",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				sellerPerp := types.NewPerpetual(0, MarketId, sellerAcc.Owner, qty10, p100, calcMargin(qty10, p100), math.LegacyZeroDec()) // Seller +10 @ 100 (M=100)
				suite.SetPerpetualStateWithEntryFR(sellerPerp)
				market.TotalOpen = qty10
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty15, p105, b, s, false) }, // B buys 15, S sells 15 @ 105
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: true, // Seller flipped to -5
				BuyerFinalQty: qty15, BuyerFinalEP: p105, BuyerFinalMargin: calcMargin(qty15, p105), // Q=15, EP=105, M=157.5
				SellerFinalQty: qty5.Neg(), SellerFinalEP: p105, SellerFinalMargin: calcMargin(qty5, p105), // Q=-5, EP=105, M=53
				BuyerBalanceChange: math.LegacyMustNewDecFromStr("-157.5").MulInt(denomFactor).TruncateInt(), // -Margin
				// Seller: Close +10: PNL=+10*(105-100)=+50. Refund=100. Open -5: Margin=-52 (after rounding). Total = +50+100-52 = +98
				SellerBalanceChange: math.LegacyMustNewDecFromStr("97.5").MulInt(denomFactor).TruncateInt(),
				MarketBalanceChange: math.NewInt(60).Mul(denomFactor), // +B_M(158) - S_PNL(50) - S_Refund(100) + S_NewM(53) = 60
				FinalMarketOI:       qty10,                            // Unchanged (1 Open, 1 Flip) = 15
				TwapBlockDataSet:    true, BuyerFundingRateSet: true, SellerFundingRateSet: true,
			},
		},
		// B(0) S(-) -> B Open L, S Inc S
		{
			name: "B(0), S(-): Buyer Opens Long, Seller Increases Short",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				sellerPerp := types.NewPerpetual(0, MarketId, sellerAcc.Owner, qty10.Neg(), p100, calcMargin(qty10, p100), math.LegacyZeroDec()) // Seller -10 @ 100 (M=100)
				suite.SetPerpetualStateWithEntryFR(sellerPerp)
				market.TotalOpen = qty10
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty5, p95, b, s, false) }, // B buys 5, S sells 5 @ 95
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: true,
				BuyerFinalQty: qty5, BuyerFinalEP: p95, BuyerFinalMargin: calcMargin(qty5, p95), // Q=5, EP=95, M=48
				// Seller: Q=-15, EP=(-10*100+-5*95)/-15 = 98.33.., M=abs(-15)*98.33*0.1 = 147.5->148
				SellerFinalQty: qty15.Neg(), SellerFinalEP: math.LegacyMustNewDecFromStr("98.333333333333333333"), SellerFinalMargin: math.LegacyMustNewDecFromStr("147.5").MulInt(denomFactor).TruncateInt(),
				BuyerBalanceChange:  math.LegacyMustNewDecFromStr("-47.5").MulInt(denomFactor).TruncateInt(), // -Margin
				SellerBalanceChange: math.LegacyMustNewDecFromStr("-47.5").MulInt(denomFactor).TruncateInt(), // -MarginDiff = -(148-100) = -48
				MarketBalanceChange: math.NewInt(95).Mul(denomFactor),                                        // +B_M(48) + S_M_Diff(48) = 96
				FinalMarketOI:       qty15,                                                                   // OI Increases by 5 to 15
				TwapBlockDataSet:    true, BuyerFundingRateSet: true, SellerFundingRateSet: true,
			},
		},

		// === GROUP 3: Buyer Has Position, Seller Opens ===
		// B(+) S(0) -> B Inc L, S Open S
		{
			name: "B(+), S(0): Buyer Increases Long, Seller Opens Short",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				buyerPerp := types.NewPerpetual(0, MarketId, buyerAcc.Owner, qty10, p100, calcMargin(qty10, p100), math.LegacyZeroDec()) // Buyer +10 @ 100 (M=100)
				suite.SetPerpetualStateWithEntryFR(buyerPerp)
				market.TotalOpen = qty10
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty5, p105, b, s, false) }, // B buys 5, S sells 5 @ 105
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: true,
				// Buyer: Q=15, EP=(10*100+5*105)/15=101.66.., M=15*101.66*0.1=152.5->153. MarginDiff=53
				BuyerFinalQty: qty15, BuyerFinalEP: math.LegacyMustNewDecFromStr("101.666666666666666667"), BuyerFinalMargin: math.LegacyMustNewDecFromStr("152.5").MulInt(denomFactor).TruncateInt(),
				SellerFinalQty: qty5.Neg(), SellerFinalEP: p105, SellerFinalMargin: calcMargin(qty5, p105), // Q=-5, EP=105, M=53
				BuyerBalanceChange:  math.LegacyMustNewDecFromStr("-52.5").MulInt(denomFactor).TruncateInt(), // -MarginDiff
				SellerBalanceChange: math.LegacyMustNewDecFromStr("-52.5").MulInt(denomFactor).TruncateInt(), // -Margin
				MarketBalanceChange: math.NewInt(105).Mul(denomFactor),                                       // +B_M_Diff(52) + S_M(52) = 104
				FinalMarketOI:       qty15,                                                                   // OI Increases by 5 to 15
				TwapBlockDataSet:    true, BuyerFundingRateSet: true, SellerFundingRateSet: true,
			},
		},
		// B(-) S(0) -> B Dec S, S Open S
		{
			name: "B(-), S(0): Buyer Decreases Short, Seller Opens Short",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				buyerPerp := types.NewPerpetual(0, MarketId, buyerAcc.Owner, qty15.Neg(), p100, calcMargin(qty15, p100), math.LegacyZeroDec()) // Buyer -15 @ 100 (M=150)
				suite.SetPerpetualStateWithEntryFR(buyerPerp)
				market.TotalOpen = qty15
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty10, p95, b, s, false) }, // B buys 10, S sells 10 @ 95
			expected: ExpectedState{
				Err: "", BuyerExists: true, SellerExists: true, // Buyer remains -5
				BuyerFinalQty: qty5.Neg(), BuyerFinalEP: p100, BuyerFinalMargin: calcMargin(qty5, p100), // Q=-5, EP=100, M=50
				SellerFinalQty: qty10.Neg(), SellerFinalEP: p95, SellerFinalMargin: calcMargin(qty10, p95), // Q=-10, EP=95, M=95
				// Buyer: PNL=-10*(95-100)=+50. Margin Refund = 150-50 = +100. Total=+150
				BuyerBalanceChange:  math.NewInt(150).Mul(denomFactor),
				SellerBalanceChange: math.NewInt(-95).Mul(denomFactor), // -Margin
				MarketBalanceChange: math.NewInt(-55).Mul(denomFactor), // -B_PNL(50) - B_Refund(100) + S_M(95) = -55
				FinalMarketOI:       qty15,                             // Unchanged (1 Dec, 1 Open) = 15
				TwapBlockDataSet:    true, BuyerFundingRateSet: true, SellerFundingRateSet: true,
			},
		},
		// Add B Close S / S Open S; B Flip S->L / S Open S cases...

		// === GROUP 4: Both Have Positions ===
		// Add more transfer types, e.g., B(-) S(+) -> B Dec S / S Dec L (Both Decrease) covered
		// B(-) S(+) -> B Close S / S Close L (Both Close) covered
		// B(+) S(+) -> B Inc L / S Dec L (Transfer Long) covered
		// B(-) S(-) -> B Dec S / S Inc S (Transfer Short) covered
		{
			name: "Transfer: Buyer Closes Long, Seller Increases Long", // Less common but possible
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				buyerPerp := types.NewPerpetual(0, MarketId, buyerAcc.Owner, qty10, p100, calcMargin(qty10, p100), math.LegacyZeroDec()) // Buyer +10 @ 100 (M=100)
				sellerPerp := types.NewPerpetual(0, MarketId, sellerAcc.Owner, qty5, p100, calcMargin(qty5, p100), math.LegacyZeroDec()) // Seller +5 @ 100 (M=50)
				suite.SetPerpetualStateWithEntryFR(buyerPerp)
				suite.SetPerpetualStateWithEntryFR(sellerPerp)
				market.TotalOpen = qty15
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)
				return buyerAcc, sellerAcc
			},
			// Seller must be buyer in trade to increase long
			trade: func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty10, p105, s, b, false) }, // S buys 10, B sells 10 @ 105
			expected: ExpectedState{
				Err: "", BuyerExists: false, SellerExists: true, // Buyer closed, Seller increased
				// Buyer: PNL=+10*(105-100)=+50. Refund=100. Total=+150
				BuyerBalanceChange: math.NewInt(150).Mul(denomFactor),
				// Seller: Q=15, EP=(5*100+10*105)/15=103.33.., M=15*103.33*0.1=155. MarginDiff=105
				SellerFinalQty: qty15, SellerFinalEP: math.LegacyMustNewDecFromStr("103.333333333333333333"), SellerFinalMargin: math.NewInt(155).Mul(denomFactor),
				SellerBalanceChange: math.NewInt(-105).Mul(denomFactor), // -MarginDiff
				MarketBalanceChange: math.NewInt(-45).Mul(denomFactor),  // -B_PNL(50) - B_Refund(100) + S_M_Diff(105) = -45
				FinalMarketOI:       qty15,                              // Unchanged (1 Close, 1 Inc) = 15
				TwapBlockDataSet:    true, BuyerFundingRateSet: false, SellerFundingRateSet: true,
			},
		},
		// Add Both Flip cases...

		// === GROUP 5: Error Cases ===
		// Keep existing error cases for funds, margin checks
		{
			name: "Error: Market Not Found",
			setup: func() (types.SubAccount, types.SubAccount) {
				_, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				buyerAcc.MarketId = 2
				sellerAcc.MarketId = 2
				return buyerAcc, sellerAcc // Return dummy market for trade obj
			},
			trade:    func(b, s types.SubAccount) types.Trade { return types.NewTrade(2, qty10, p100, b, s, false) },
			expected: ExpectedState{Err: types.ErrPerpetualMarketNotFound.Error()}, // Use specific error
		},
		{
			name: "Error: Trade and sub account market mismatch",
			setup: func() (types.SubAccount, types.SubAccount) {
				_, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				return buyerAcc, sellerAcc // Return dummy market for trade obj
			},
			trade:    func(b, s types.SubAccount) types.Trade { return types.NewTrade(2, qty10, p100, b, s, false) },
			expected: ExpectedState{Err: "trade market id and subAccounts market id does not match"}, // Use specific error
		},
		{
			name: "Error: Trade quantity <= 0",
			setup: func() (types.SubAccount, types.SubAccount) {
				_, buyerAcc, sellerAcc, _ := suite.SetupExchangeTest()
				return buyerAcc, sellerAcc // Return dummy market for trade obj
			},
			trade: func(b, s types.SubAccount) types.Trade {
				return types.Trade{b, s, MarketId, p105, math.LegacyZeroDec(), false}
			},
			expected: ExpectedState{Err: "trade quantity must be greater than zero"}, // Use specific error
		},
		{
			name: "SettleFunding Fails - Market Insufficient Funds to Pay Long",
			setup: func() (types.SubAccount, types.SubAccount) {
				market, buyerAcc, sellerAcc, marketAccAddr := suite.SetupExchangeTest()

				prevBlockHeight := suite.ctx.BlockHeight()
				prevBlockTime := suite.ctx.BlockTime().Add(-10 * time.Second) // Example time step
				prevCtx := suite.ctx.WithBlockHeight(prevBlockHeight).WithBlockTime(prevBlockTime)

				dummyTrade := types.NewTrade(MarketId, math.LegacyNewDec(1), math.LegacyNewDec(100), buyerAcc, sellerAcc, false)
				err := suite.app.ClobKeeper.SetTwapPrices(prevCtx, dummyTrade)
				suite.Require().NoError(err, "Setup: Failed to set previous TWAP price")

				suite.IncreaseHeight(1)
				err = suite.app.ClobKeeper.SetTwapPrices(suite.ctx, dummyTrade)
				suite.Require().NoError(err, "Setup: Failed to set previous TWAP price")

				currentBlockHeight := suite.ctx.BlockHeight()

				currentFundingRate := types.FundingRate{MarketId: MarketId, Rate: math.LegacyZeroDec(), Block: uint64(currentBlockHeight)}
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, currentFundingRate)

				buyerPerp := types.NewPerpetual(0, MarketId, buyerAcc.Owner, qty10, p100, calcMargin(qty10, p100), math.LegacyZeroDec())
				suite.SetPerpetualStateWithEntryFR(buyerPerp)
				perpState, found := suite.GetPerpetualState(buyerAcc.GetOwnerAccAddress(), MarketId)
				suite.Require().True(found, "Perpetual state not found")
				// Setup funding rate difference to cause payout FROM market
				oldRate := math.LegacyZeroDec()
				perpState.EntryFundingRate = oldRate
				suite.app.ClobKeeper.SetPerpetual(suite.ctx, perpState)
				perpState, found = suite.GetPerpetualState(buyerAcc.GetOwnerAccAddress(), MarketId)
				suite.Require().True(found, "Perpetual state not found")
				suite.Require().Equal(perpState.EntryFundingRate, oldRate, "Perpetual state entry funding rate mismatch")
				market.TotalOpen = qty10
				suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)                                                                                                                 // Set buyer state with old rate
				suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{MarketId: MarketId, Rate: math.LegacyMustNewDecFromStr("-0.01"), Block: uint64(suite.ctx.BlockHeight())}) // Set current rate to 0
				// Empty the market account
				err = suite.BurnAccountBalance(marketAccAddr, QuoteDenom)
				suite.Require().NoError(err, "Failed to empty market account")
				suite.Require().True(
					suite.GetAccountBalance(marketAccAddr, QuoteDenom).IsZero(),
					"Market account balance MUST be zero after burn for this test case",
				)
				return buyerAcc, sellerAcc
			},
			trade:    func(b, s types.SubAccount) types.Trade { return types.NewTrade(MarketId, qty5, p105, b, s, false) }, // Buyer increases long, triggers SettleFunding first
			expected: ExpectedState{Err: "insufficient funds"},                                                             // Error expected from AddToSubAccount inside SettleFunding
		},
		// Add GetPerpetual error case (might require mocking or complex setup)
		// Add SetTwapPrices error case (might require mocking market lookup inside it)
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// --- Setup ---
			buyerAcc, sellerAcc := tc.setup()
			startCtx := suite.ctx // Capture context after setup
			currentBlockHeight := uint64(startCtx.BlockHeight())
			currentFundingRate := suite.app.ClobKeeper.GetFundingRate(startCtx, MarketId)

			market, err := suite.app.ClobKeeper.GetPerpetualMarket(startCtx, MarketId)
			suite.Require().NoError(err, "Market not found after setup for case: %s", tc.name)
			initialMarketOI := market.TotalOpen // THIS is the correct initial OI

			initialBuyerBalance := suite.GetAccountBalance(buyerAcc.GetTradingAccountAddress(), QuoteDenom)
			initialSellerBalance := suite.GetAccountBalance(sellerAcc.GetTradingAccountAddress(), QuoteDenom)
			initialMarketBalance := suite.GetAccountBalance(market.GetAccount(), QuoteDenom)

			trade := tc.trade(buyerAcc, sellerAcc)

			// --- Execute ---
			err = suite.app.ClobKeeper.Exchange(startCtx, trade)

			// --- Assert ---
			if tc.expected.Err != "" {
				suite.Require().Error(err, "Expected an error but got nil")
				suite.Require().Contains(err.Error(), tc.expected.Err, "Error message mismatch")

				// Verify balances did not change on error *relative to start of this test case*
				suite.CheckBalanceChange(buyerAcc.GetTradingAccountAddress(), initialBuyerBalance, math.ZeroInt(), "Buyer (on error)")
				suite.CheckBalanceChange(sellerAcc.GetTradingAccountAddress(), initialSellerBalance, math.ZeroInt(), "Seller (on error)")
				suite.CheckBalanceChange(market.GetAccount(), initialMarketBalance, math.ZeroInt(), "Market (on error)")

				// Check OI unchanged
				finalMarket, err := suite.app.ClobKeeper.GetPerpetualMarket(startCtx, MarketId)
				if err == nil { // Market might not be found if that was the error
					suite.Require().True(initialMarketOI.Equal(finalMarket.TotalOpen), "Market OI changed despite error")
				}
				// Potentially check perpetual states haven't changed from setup state either

			} else {
				suite.Require().NoError(err, "Expected no error but got: %v", err)

				// Check Buyer Final State
				finalBuyerPerp, buyerFound := suite.GetPerpetualState(buyerAcc.GetOwnerAccAddress(), MarketId)
				suite.Require().Equal(tc.expected.BuyerExists, buyerFound, "Buyer existence mismatch")
				if tc.expected.BuyerExists {
					suite.Require().True(tc.expected.BuyerFinalQty.Equal(finalBuyerPerp.Quantity), "Buyer Qty mismatch: Exp %s Got %s", tc.expected.BuyerFinalQty, finalBuyerPerp.Quantity)
					suite.Require().True(tc.expected.BuyerFinalEP.Equal(finalBuyerPerp.EntryPrice), "Buyer EP mismatch: Exp %s Got %s", tc.expected.BuyerFinalEP, finalBuyerPerp.EntryPrice)
					suite.Require().True(tc.expected.BuyerFinalMargin.Equal(finalBuyerPerp.MarginAmount), "Buyer Margin mismatch: Exp %s Got %s", tc.expected.BuyerFinalMargin, finalBuyerPerp.MarginAmount)
					if tc.expected.BuyerFundingRateSet {
						suite.Require().True(currentFundingRate.Rate.Equal(finalBuyerPerp.EntryFundingRate), "Buyer EntryFundingRate mismatch: Exp %s Got %s", currentFundingRate.Rate, finalBuyerPerp.EntryFundingRate)
					}
				} else {
					_, ownerFound := suite.app.ClobKeeper.GetPerpetualOwner(startCtx, buyerAcc.GetOwnerAccAddress(), MarketId)
					suite.Require().False(ownerFound, "Buyer owner mapping should be deleted but was found")
				}

				// Check Seller Final State
				finalSellerPerp, sellerFound := suite.GetPerpetualState(sellerAcc.GetOwnerAccAddress(), MarketId)
				suite.Require().Equal(tc.expected.SellerExists, sellerFound, "Seller existence mismatch")
				if tc.expected.SellerExists {
					suite.Require().True(tc.expected.SellerFinalQty.Equal(finalSellerPerp.Quantity), "Seller Qty mismatch: Exp %s Got %s", tc.expected.SellerFinalQty, finalSellerPerp.Quantity)
					suite.Require().True(tc.expected.SellerFinalEP.Equal(finalSellerPerp.EntryPrice), "Seller EP mismatch: Exp %s Got %s", tc.expected.SellerFinalEP, finalSellerPerp.EntryPrice)
					suite.Require().True(tc.expected.SellerFinalMargin.Equal(finalSellerPerp.MarginAmount), "Seller Margin mismatch: Exp %s Got %s", tc.expected.SellerFinalMargin, finalSellerPerp.MarginAmount)
					if tc.expected.SellerFundingRateSet {
						suite.Require().True(currentFundingRate.Rate.Equal(finalSellerPerp.EntryFundingRate), "Seller EntryFundingRate mismatch: Exp %s Got %s", currentFundingRate.Rate, finalSellerPerp.EntryFundingRate)
					}
				} else {
					_, ownerFound := suite.app.ClobKeeper.GetPerpetualOwner(startCtx, sellerAcc.GetOwnerAccAddress(), MarketId)
					suite.Require().False(ownerFound, "Seller owner mapping should be deleted but was found")
				}

				// Check Balances
				suite.CheckBalanceChange(buyerAcc.GetTradingAccountAddress(), initialBuyerBalance, tc.expected.BuyerBalanceChange, "Buyer")
				suite.CheckBalanceChange(sellerAcc.GetTradingAccountAddress(), initialSellerBalance, tc.expected.SellerBalanceChange, "Seller")
				suite.CheckBalanceChange(market.GetAccount(), initialMarketBalance, tc.expected.MarketBalanceChange, "Market")

				// Check Market OI
				finalMarket, err := suite.app.ClobKeeper.GetPerpetualMarket(startCtx, MarketId)
				suite.Require().NoError(err, "Market should be found after successful execution")
				// Assert final OI against the expected value (which was calculated based on the true initial OI)
				suite.Require().True(tc.expected.FinalMarketOI.Equal(finalMarket.TotalOpen), "Market OI mismatch. Expected %s, Got %s (True Initial was %s)", tc.expected.FinalMarketOI, finalMarket.TotalOpen, initialMarketOI)

				// Check TWAP
				if tc.expected.TwapBlockDataSet {
					allTwap := suite.app.ClobKeeper.GetAllTwapPrices(startCtx) // Use context after execution
					foundTwapForBlock := false
					for _, twapPrice := range allTwap {
						// Check if a record exists for the correct market and block height
						if twapPrice.MarketId == MarketId && twapPrice.Block == currentBlockHeight {
							foundTwapForBlock = true
							break // Found it, no need to check further
						}
					}
					// Assert that a record for the executed block height was found
					suite.Require().True(foundTwapForBlock, "Expected TWAP data for block %d but none found in GetAllTwapPrices result", currentBlockHeight)
				}
			}
		})
	}
}
