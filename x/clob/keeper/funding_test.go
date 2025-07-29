package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (suite *KeeperTestSuite) TestSettleFunding() {
	// --- Common Setup ---
	initialEntryRate := math.LegacyMustNewDecFromStr("0.01")     // 1%
	currentRatePosDelta := math.LegacyMustNewDecFromStr("0.03")  // 3% (Delta = +0.02)
	currentRateNegDelta := math.LegacyMustNewDecFromStr("-0.01") // -1% (Delta = -0.02)
	currentRateZeroDelta := initialEntryRate                     // 1% (Delta = 0)

	baseTwapPrice := math.LegacyNewDec(100)
	posQty := math.LegacyNewDec(10)
	notionalValue := posQty.Mul(baseTwapPrice) // 10 * 100 = 1000

	// Calculate expected payment magnitude based on rate deltas
	posRateDelta := currentRatePosDelta.Sub(initialEntryRate)           // +0.02
	negRateDelta := currentRateNegDelta.Sub(initialEntryRate)           // -0.02
	paymentAmountPosDelta := posRateDelta.Mul(notionalValue).RoundInt() // +0.02 * 1000 = 20
	paymentAmountNegDelta := negRateDelta.Mul(notionalValue).RoundInt() // -0.02 * 1000 = -20 -> Abs() = 20

	// Define Perpetual Templates
	margin := math.NewInt(100) // Example margin
	longPerpBase := types.NewPerpetual(1, MarketId, "owner_placeholder", posQty, math.LegacyNewDec(99), margin, initialEntryRate, MarketId)
	shortPerpBase := types.NewPerpetual(2, MarketId, "owner_placeholder", posQty.Neg(), math.LegacyNewDec(101), margin, initialEntryRate, MarketId)

	testCases := []struct {
		name                        string
		currentRate                 types.FundingRate                                             // Rate set in store for GetFundingRate to return
		perpetual                   types.Perpetual                                               // Initial perpetual state (copy passed)
		mockTwapPrice               math.LegacyDec                                                // Value GetCurrentTwapPrice should return
		initialSubAccBal            math.Int                                                      // Needed for setup/verification
		initialMarketBal            math.Int                                                      // Needed for setup/verification
		expectedFinalEntryRate      math.LegacyDec                                                // Expected perpetual.EntryFundingRate after call
		expectedSubAccBalanceChange math.Int                                                      // Net change for the subAccount
		expectedMarketBalanceChange math.Int                                                      // Net change for the market account
		expectedErr                 string                                                        // Error substring or ""
		setupFunc                   func(subAccAddr sdk.AccAddress, marketAccAddr sdk.AccAddress) // Optional setup specific to the test case
	}{
		// --- Long Position Cases ---
		{
			name:             "Long receives funding (Negative Rate Delta)",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRateNegDelta, Block: 2},
			perpetual:        longPerpBase, // Entry=0.01
			mockTwapPrice:    baseTwapPrice,
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate: currentRateNegDelta,
			// Delta = -0.02. fundingPnL = (-0.02*1000)*(-1) = +20. AddToSubAccount called.
			expectedSubAccBalanceChange: paymentAmountPosDelta,       // +20
			expectedMarketBalanceChange: paymentAmountPosDelta.Neg(), // -20
			expectedErr:                 "",
		},
		{
			name:             "Long pays funding (Positive Rate Delta)",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRatePosDelta, Block: 2},
			perpetual:        longPerpBase, // Entry=0.01
			mockTwapPrice:    baseTwapPrice,
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate: currentRatePosDelta,
			// Delta = +0.02. fundingPnL = (+0.02*1000)*(-1) = -20. SendFromSubAccount called.
			expectedSubAccBalanceChange: paymentAmountNegDelta,       // -20
			expectedMarketBalanceChange: paymentAmountNegDelta.Neg(), // +20
			expectedErr:                 "",
		},
		{
			name:             "Long zero funding (Zero Rate Delta)",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRateZeroDelta, Block: 2},
			perpetual:        longPerpBase, // Entry=0.01
			mockTwapPrice:    baseTwapPrice,
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate: currentRateZeroDelta,
			// Delta = 0.01 - 0.01 = 0. PNL = 0. No transfer.
			expectedSubAccBalanceChange: math.ZeroInt(),
			expectedMarketBalanceChange: math.ZeroInt(),
			expectedErr:                 "",
		},

		// --- Short Position Cases ---
		{
			name:             "Short pays funding (Negative Rate Delta)",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRateNegDelta, Block: 2},
			perpetual:        shortPerpBase, // Entry=0.01
			mockTwapPrice:    baseTwapPrice,
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate: currentRateNegDelta,
			// Delta = -0.02. PNL = (-0.02*1000)*(+1) = -20. SendFromSubAccount called.
			expectedSubAccBalanceChange: paymentAmountNegDelta,       // -20
			expectedMarketBalanceChange: paymentAmountNegDelta.Neg(), // +20
			expectedErr:                 "",
		},
		{
			name:             "Short receives funding (Positive Rate Delta)",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRatePosDelta, Block: 2},
			perpetual:        shortPerpBase, // Entry=0.01
			mockTwapPrice:    baseTwapPrice,
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate: currentRatePosDelta,
			// Delta = +0.02. PNL = (+0.02*1000)*(+1) = +20. AddToSubAccount called.
			expectedSubAccBalanceChange: paymentAmountPosDelta,       // +20
			expectedMarketBalanceChange: paymentAmountPosDelta.Neg(), // -20
			expectedErr:                 "",
		},
		{
			name:             "Short zero funding (Zero Rate Delta)",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRateZeroDelta, Block: 2},
			perpetual:        shortPerpBase, // Entry=0.01
			mockTwapPrice:    baseTwapPrice,
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate:      currentRateZeroDelta,
			expectedSubAccBalanceChange: math.ZeroInt(),
			expectedMarketBalanceChange: math.ZeroInt(),
			expectedErr:                 "",
		},

		// --- Edge Cases ---
		{
			name:             "Zero TWAP Price results in zero funding",
			currentRate:      types.FundingRate{MarketId: MarketId, Rate: currentRatePosDelta, Block: 2}, // Rate delta is non-zero (+0.02)
			perpetual:        longPerpBase,
			mockTwapPrice:    math.LegacyZeroDec(), // Force TWAP to zero
			initialSubAccBal: math.NewInt(100000), initialMarketBal: math.NewInt(100000),
			expectedFinalEntryRate:      currentRatePosDelta, // Entry rate still updates
			expectedSubAccBalanceChange: math.ZeroInt(),      // PNL calculation becomes zero
			expectedMarketBalanceChange: math.ZeroInt(),
			expectedErr:                 "",
		},

		// --- Error Cases ---
		{
			name:                        "Error: SubAccount insufficient funds to pay funding",
			currentRate:                 types.FundingRate{MarketId: MarketId, Rate: currentRatePosDelta, Block: 2}, // Long Pays 20
			perpetual:                   longPerpBase,
			mockTwapPrice:               baseTwapPrice,
			initialSubAccBal:            math.NewInt(10), // Less than payment of 20
			initialMarketBal:            math.NewInt(100000),
			expectedFinalEntryRate:      initialEntryRate,     // Rate not updated on error
			expectedSubAccBalanceChange: math.ZeroInt(),       // Balance unchanged
			expectedMarketBalanceChange: math.ZeroInt(),       // Balance unchanged
			expectedErr:                 "insufficient funds", // From SendFromSubAccount
			setupFunc: func(subAccAddr sdk.AccAddress, marketAccAddr sdk.AccAddress) {
				// Set specific low balance for sub account
				suite.SetAccountBalance(subAccAddr, sdk.NewCoins(sdk.NewCoin(QuoteDenom, math.NewInt(10))))
			},
		},
		{
			// This test's success depends on AddToSubAccount checking source balance.
			// If AddToSubAccount mints, this test should be removed or expect success.
			name:                        "Error: Market insufficient funds to pay funding (if Add sends)",
			currentRate:                 types.FundingRate{MarketId: MarketId, Rate: currentRateNegDelta, Block: 2}, // Long Receives 20
			perpetual:                   longPerpBase,
			mockTwapPrice:               baseTwapPrice,
			initialSubAccBal:            math.NewInt(100000),
			initialMarketBal:            math.NewInt(10),      // Less than payment of 20
			expectedFinalEntryRate:      initialEntryRate,     // Rate not updated on error
			expectedSubAccBalanceChange: math.ZeroInt(),       // Balance unchanged
			expectedMarketBalanceChange: math.ZeroInt(),       // Balance unchanged
			expectedErr:                 "insufficient funds", // From AddToSubAccount (if it sends from marketAcc)
			setupFunc: func(subAccAddr sdk.AccAddress, marketAccAddr sdk.AccAddress) {
				// Set specific low balance for market account
				suite.SetAccountBalance(marketAccAddr, sdk.NewCoins(sdk.NewCoin(QuoteDenom, math.NewInt(10))))
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			// --- Setup ---
			// Setup requires a valid market and subaccount
			market, subAcc, _, marketAccAddr := suite.SetupExchangeTest() // Use helper to get accounts/market

			// Set initial balances for subaccount and market account for this test case run
			suite.SetAccountBalance(subAcc.GetTradingAccountAddress(), sdk.NewCoins(sdk.NewCoin(QuoteDenom, tc.initialSubAccBal)))
			suite.SetAccountBalance(marketAccAddr, sdk.NewCoins(sdk.NewCoin(QuoteDenom, tc.initialMarketBal)))

			// Run specific setup if provided (e.g., adjusting balances further)
			if tc.setupFunc != nil {
				tc.setupFunc(subAcc.GetTradingAccountAddress(), marketAccAddr)
			}

			// Set the current funding rate in the store for the keeper to read
			suite.app.ClobKeeper.SetFundingRate(suite.ctx, tc.currentRate)

			// Make a copy of the perpetual to pass its pointer, and ensure owner matches subAcc
			perpetualCopy := tc.perpetual
			perpetualCopy.Owner = subAcc.Owner                         // Assign correct owner
			initialEntryRateForCheck := perpetualCopy.EntryFundingRate // Store before call

			// Capture initial balances *after* all setup adjustments
			initialSubAccBalanceCheck := suite.GetAccountBalance(subAcc.GetTradingAccountAddress(), QuoteDenom)
			initialMarketBalanceCheck := suite.GetAccountBalance(marketAccAddr, QuoteDenom)

			suite.SetTwapPriceDirectly(MarketId, tc.mockTwapPrice)

			// --- Execute ---
			// Pass pointer to the subAccount struct (though Send/Add use by value internally based on prev code)
			// Pass pointer to perpetualCopy struct
			err := suite.app.ClobKeeper.SettleFunding(suite.ctx, &subAcc, market, &perpetualCopy)

			// --- Assert ---
			if tc.expectedErr != "" {
				suite.Require().Error(err, "Expected an error but got nil")
				suite.Require().Contains(err.Error(), tc.expectedErr, "Error message mismatch")

				// Check state was NOT updated on error
				suite.Require().True(initialEntryRateForCheck.Equal(perpetualCopy.EntryFundingRate), "EntryFundingRate changed despite error")
				// Check balances unchanged from the state *just before* the call
				suite.CheckBalanceChange(subAcc.GetTradingAccountAddress(), initialSubAccBalanceCheck, math.ZeroInt(), "SubAccount (on error)")
				suite.CheckBalanceChange(marketAccAddr, initialMarketBalanceCheck, math.ZeroInt(), "Market (on error)")

			} else {
				suite.Require().NoError(err, "Expected no error but got: %v", err)

				// Check EntryFundingRate was updated correctly on the copy (caller needs to persist)
				suite.Require().True(tc.expectedFinalEntryRate.Equal(perpetualCopy.EntryFundingRate), "Final EntryFundingRate mismatch. Expected %s, Got %s", tc.expectedFinalEntryRate, perpetualCopy.EntryFundingRate)

				// Check Balances changed correctly from the state *just before* the call
				suite.CheckBalanceChange(subAcc.GetTradingAccountAddress(), initialSubAccBalanceCheck, tc.expectedSubAccBalanceChange, "SubAccount")
				suite.CheckBalanceChange(marketAccAddr, initialMarketBalanceCheck, tc.expectedMarketBalanceChange, "Market")

				// Optional: Verify persisted state if desired
				suite.Require().True(tc.expectedFinalEntryRate.Equal(perpetualCopy.EntryFundingRate), "Persisted EntryFundingRate mismatch")
			}
			// suite.mockTwapEnabled = false // Disable mock if using that pattern
		})
	}
}
