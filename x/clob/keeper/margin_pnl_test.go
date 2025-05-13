package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

var (
	IMR = math.LegacyMustNewDecFromStr("0.1")
)

func (suite *KeeperTestSuite) TestSettleMarginAndRPnL_Complete() { // Renamed slightly
	suite.ResetSuite()

	// --- Setup ---
	markets := suite.CreateMarket("uatom")
	market := markets[0]
	marketAccAddr := market.GetAccount() // Get market module account address

	initialBalanceAmt := int64(100_000_000)       // Generous balance for traders
	initialMarketBalanceAmt := int64(500_000_000) // Ensure market has enough funds to pay out PNL/Margin
	initialBalance := sdk.NewCoin(QuoteDenom, math.NewInt(initialBalanceAmt))
	initialMarketBalanceCoin := sdk.NewCoin(QuoteDenom, math.NewInt(initialMarketBalanceAmt))

	subAccounts := suite.SetupSubAccounts(2, sdk.NewCoins(initialBalance))
	buyerAcc := subAccounts[0]
	sellerAcc := subAccounts[1]

	// Fund market account adequately - CRITICAL FIX
	suite.FundAccount(marketAccAddr, sdk.NewCoins(initialMarketBalanceCoin))

	// --- Test Cases ---
	testCases := []struct {
		name                      string
		oldPerpetual              types.Perpetual // Use Zero value for new position
		trade                     types.Trade
		isBuyer                   bool
		expectedUpdatedQuantity   math.LegacyDec
		expectedUpdatedEntryPrice math.LegacyDec
		expectedUpdatedMargin     math.Int // Changed to math.Int
		// For balance checks:
		expectedSubAccBalanceChange math.Int // Net change (PNL - MarginDelta or PNL + MarginDelta) - Changed to math.Int
		expectedMarketBalanceChange math.Int // Should be opposite of SubAcc change - Changed to math.Int
		expectedErr                 bool
		expectedErrMsg              string
	}{
		// --- Case 1: Increase Short ---
		{
			name:                        "Case 1: Seller increases Short position",
			oldPerpetual:                types.NewPerpetual(1, MarketId, sellerAcc.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000), math.LegacyZeroDec()),
			trade:                       types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(90), buyerAcc, sellerAcc),
			isBuyer:                     false,
			expectedUpdatedQuantity:     math.LegacyNewDec(-15),
			expectedUpdatedEntryPrice:   math.LegacyMustNewDecFromStr("96.666666666666666667"), // Avg EP = (-10*100 + -5*90)/-15 = 96.66...
			expectedUpdatedMargin:       math.NewInt(145_000_000),                              // New Margin = 15 * 96.66... * 0.1 = 145
			expectedSubAccBalanceChange: math.NewInt(-45_000_000),                              // PNL=0, MarginDelta = -(145-100) = -45
			expectedMarketBalanceChange: math.NewInt(45_000_000),                               // Receives margin diff
			expectedErr:                 false,
		},

		// --- Case 2: Decrease Short ---
		{
			name:                      "Case 2: Buyer decreases Short position",
			oldPerpetual:              types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(-20), math.LegacyNewDec(100), math.NewInt(200_000_000), math.LegacyZeroDec()),
			trade:                     types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(90), buyerAcc, sellerAcc),
			isBuyer:                   true,
			expectedUpdatedQuantity:   math.LegacyNewDec(-15),
			expectedUpdatedEntryPrice: math.LegacyNewDec(100),
			expectedUpdatedMargin:     math.NewInt(150_000_000), // Remaining Margin = 15 * 100 * 0.1 = 150
			// PNL = -5 * (90 - 100) = +50
			// MarginDelta = +(200 - 150) = +50 (Refund)
			expectedSubAccBalanceChange: math.NewInt(100_000_000),  // PNL + MarginDelta = 50 + 50
			expectedMarketBalanceChange: math.NewInt(-100_000_000), // Pays PNL and Margin Refund
			expectedErr:                 false,
		},

		// --- Case 3: Close Short ---
		{
			name:                      "Case 3: Buyer closes Short position",
			oldPerpetual:              types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000), math.LegacyZeroDec()),
			trade:                     types.NewTrade(MarketId, math.LegacyNewDec(10), math.LegacyNewDec(95), buyerAcc, sellerAcc), // Buy 10 @ 95
			isBuyer:                   true,
			expectedUpdatedQuantity:   math.LegacyZeroDec(),
			expectedUpdatedEntryPrice: math.LegacyZeroDec(),
			expectedUpdatedMargin:     math.ZeroInt(),
			// PNL = -10 * (95 - 100) = +50
			// MarginDelta = +100 (Full Refund)
			expectedSubAccBalanceChange: math.NewInt(150_000_000),  // PNL + MarginDelta = 50 + 100
			expectedMarketBalanceChange: math.NewInt(-150_000_000), // Pays PNL and Margin Refund
			expectedErr:                 false,
		},

		// --- Case 4: Flip Short -> Long ---
		{
			name:                      "Case 4: Buyer flips Short to Long",
			oldPerpetual:              types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000), math.LegacyZeroDec()),
			trade:                     types.NewTrade(MarketId, math.LegacyNewDec(15), math.LegacyNewDec(98), buyerAcc, sellerAcc), // Buy 15 @ 98
			isBuyer:                   true,
			expectedUpdatedQuantity:   math.LegacyNewDec(5),
			expectedUpdatedEntryPrice: math.LegacyNewDec(98),
			expectedUpdatedMargin:     math.NewInt(49_000_000), // New Margin = 5 * 98 * 0.1 = 49
			// PNL = -10 * (98 - 100) = +20
			// MarginDelta = +100 (Refund old) - 49 (Deduct new) = +51
			expectedSubAccBalanceChange: math.NewInt(71_000_000),  // PNL + MarginDelta = 20 + 51
			expectedMarketBalanceChange: math.NewInt(-71_000_000), // Pays PNL+Refund, Receives new Margin
			expectedErr:                 false,
		},

		// --- Case 5: Increase Long ---
		{
			name:                        "Case 5: Buyer increases Long position",
			oldPerpetual:                types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000), math.LegacyZeroDec()), // Q=10, EP=100, M=100
			trade:                       types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(110), buyerAcc, sellerAcc),                                                    // Buy 5 @ 110
			isBuyer:                     true,
			expectedUpdatedQuantity:     math.LegacyNewDec(15),
			expectedUpdatedEntryPrice:   math.LegacyMustNewDecFromStr("103.333333333333333333"), // Avg EP = (10*100 + 5*110) / 15 = 103.33...
			expectedUpdatedMargin:       math.NewInt(155_000_000),                               // New Margin = 15 * 103.33... * 0.1 = 155
			expectedSubAccBalanceChange: math.NewInt(-55_000_000),                               // PNL=0, MarginDelta = -(155-100) = -55
			expectedMarketBalanceChange: math.NewInt(55_000_000),                                // Receives margin diff
			expectedErr:                 false,
		},

		// --- Case 6: Decrease Long ---
		{
			name:                      "Case 6: Seller decreases Long position",
			oldPerpetual:              types.NewPerpetual(1, MarketId, sellerAcc.Owner, math.LegacyNewDec(20), math.LegacyNewDec(100), math.NewInt(200_000_000), math.LegacyZeroDec()),
			trade:                     types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(110), buyerAcc, sellerAcc),
			isBuyer:                   false,
			expectedUpdatedQuantity:   math.LegacyNewDec(15),
			expectedUpdatedEntryPrice: math.LegacyNewDec(100),
			expectedUpdatedMargin:     math.NewInt(150_000_000), // Remaining Margin = 15 * 100 * 0.1 = 150
			// PNL = +5 * (110 - 100) = +50
			// MarginDelta = +(200 - 150) = +50 (Refund)
			expectedSubAccBalanceChange: math.NewInt(100_000_000),  // PNL + MarginDelta = 50 + 50
			expectedMarketBalanceChange: math.NewInt(-100_000_000), // Pays PNL and Margin Refund
			expectedErr:                 false,
		},

		// --- Case 7: Close Long ---
		{
			name:                      "Case 7: Seller closes Long position", // This case caused the error
			oldPerpetual:              types.NewPerpetual(1, MarketId, sellerAcc.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000), math.LegacyZeroDec()),
			trade:                     types.NewTrade(MarketId, math.LegacyNewDec(10), math.LegacyNewDec(105), buyerAcc, sellerAcc), // Sell 10 @ 105
			isBuyer:                   false,
			expectedUpdatedQuantity:   math.LegacyZeroDec(),
			expectedUpdatedEntryPrice: math.LegacyZeroDec(), // Irrelevant
			expectedUpdatedMargin:     math.ZeroInt(),       // Margin returned
			// PNL = +10 * (105 - 100) = +50
			// MarginDelta = +100 (Full Refund)
			expectedSubAccBalanceChange: math.NewInt(150_000_000),  // PNL + MarginDelta = 50 + 100
			expectedMarketBalanceChange: math.NewInt(-150_000_000), // Pays PNL (50) and Margin Refund (100)
			expectedErr:                 false,
		},

		// --- Case 8: Flip Long -> Short ---
		{
			name:                      "Case 8: Seller flips Long to Short",
			oldPerpetual:              types.NewPerpetual(1, MarketId, sellerAcc.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000), math.LegacyZeroDec()),
			trade:                     types.NewTrade(MarketId, math.LegacyNewDec(15), math.LegacyNewDec(102), buyerAcc, sellerAcc), // Sell 15 @ 102
			isBuyer:                   false,
			expectedUpdatedQuantity:   math.LegacyNewDec(-5),
			expectedUpdatedEntryPrice: math.LegacyNewDec(102),
			expectedUpdatedMargin:     math.NewInt(51_000_000), // New Margin = 5 * 102 * 0.1 = 51
			// PNL (on closing +10) = +10 * (102 - 100) = +20
			// MarginDelta = +100 (Refund old) - 51 (Deduct new) = +49
			expectedSubAccBalanceChange: math.NewInt(69_000_000),  // PNL + MarginDelta = 20 + 49
			expectedMarketBalanceChange: math.NewInt(-69_000_000), // Pays PNL+Refund, Receives new Margin
			expectedErr:                 false,
		},

		// --- Case 9: Open New Position ---
		{
			name:                        "Case 9: Buyer opens new Long position",
			oldPerpetual:                types.Perpetual{Quantity: math.LegacyZeroDec()},
			trade:                       types.NewTrade(MarketId, math.LegacyNewDec(10), math.LegacyNewDec(100), buyerAcc, sellerAcc),
			isBuyer:                     true,
			expectedUpdatedQuantity:     math.LegacyNewDec(10),
			expectedUpdatedEntryPrice:   math.LegacyNewDec(100),
			expectedUpdatedMargin:       math.NewInt(100_000_000),  // 10*100*0.1
			expectedSubAccBalanceChange: math.NewInt(-100_000_000), // PNL=0, MarginSent=-100
			expectedMarketBalanceChange: math.NewInt(100_000_000),  // Receives margin
			expectedErr:                 false,
		},
		{
			name:                        "Case 9: Seller opens new Short position",
			oldPerpetual:                types.Perpetual{Quantity: math.LegacyZeroDec()},
			trade:                       types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(200), buyerAcc, sellerAcc),
			isBuyer:                     false,
			expectedUpdatedQuantity:     math.LegacyNewDec(-5),
			expectedUpdatedEntryPrice:   math.LegacyNewDec(200),
			expectedUpdatedMargin:       math.NewInt(100_000_000),  // 5*200*0.1
			expectedSubAccBalanceChange: math.NewInt(-100_000_000), // PNL=0, MarginSent=-100
			expectedMarketBalanceChange: math.NewInt(100_000_000),  // Receives margin
			expectedErr:                 false,
		},

		// --- Error Cases ---
		{
			name:         "Error: Insufficient funds to open Long",
			oldPerpetual: types.Perpetual{Quantity: math.LegacyZeroDec()},
			// Corrected Quantity: Requires 100,000,010 margin > 100,000,000 initial balance
			trade:          types.NewTrade(MarketId, math.LegacyNewDec(10_000_001), math.LegacyNewDec(100), buyerAcc, sellerAcc),
			isBuyer:        true,
			expectedErr:    true,
			expectedErrMsg: "insufficient funds", // Or specific bank error
		},
		{
			name:         "Error: Insufficient funds to increase Long margin",
			oldPerpetual: types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(1), math.LegacyNewDec(10), math.NewInt(1_000_000), math.LegacyZeroDec()), // Start small
			// Corrected Quantity: Requires ~110,000,000 margin diff > 100,000,000 initial balance
			trade:          types.NewTrade(MarketId, math.LegacyNewDec(11_000_000), math.LegacyNewDec(100), buyerAcc, sellerAcc),
			isBuyer:        true,
			expectedErr:    true,
			expectedErrMsg: "insufficient funds", // Or specific bank error
		},
		{
			name:           "Error: Margin check fail on increase Long (new <= old)",
			oldPerpetual:   types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(1000_000_000), math.LegacyZeroDec()), // Artificially high old margin
			trade:          types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(10), buyerAcc, sellerAcc),                                                      // Low trade price
			isBuyer:        true,
			expectedErr:    true,
			expectedErrMsg: "newRequiredInitialMargin (105000000) must be greater than oldPerpetual.Margin (1000000000) for buyer when position is increased from positive to more positive", // Contains the error message from fmt.Errorf
		},
		{
			name:           "Error: Margin check fail on decrease Long (old <= new)",
			oldPerpetual:   types.NewPerpetual(1, MarketId, sellerAcc.Owner, math.LegacyNewDec(20), math.LegacyNewDec(100), math.NewInt(1_000_000), math.LegacyZeroDec()), // Artificially low old margin
			trade:          types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(110), buyerAcc, sellerAcc),
			isBuyer:        false,
			expectedErr:    true,
			expectedErrMsg: "oldPerpetual.Margin (1000000) must be greater than newRequiredInitialMargin (150000000) for seller when position is reduced from positive to less positive",
		},
		{
			name: "Error Return: Margin check fail on decrease Short (old <= new)",
			// Setup: Old margin (5) is <= required margin for remaining position (abs(-5)*100*0.1 = 50) -> Should fail
			oldPerpetual: types.NewPerpetual(1, MarketId, buyerAcc.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(5_000_000), math.LegacyZeroDec()), // Artificially low old margin
			trade:        types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(90), buyerAcc, sellerAcc),                                                    // Buy 5 @ 90 -> remaining is -5 Qty
			isBuyer:      true,
			// Expecting error return
			expectedErrMsg: "oldPerpetual.Margin (5000000) must be greater than newRequiredInitialMargin (50000000) for buyer when position is reduced from negative to less negative", // Match fmt.Errorf string
		},
		{
			name: "Error Return: Margin check fail on increase Short (new <= old)",
			// Setup: Required margin for new pos (abs(-15)*70*0.1=105) is <= old margin (1000) -> Should fail
			oldPerpetual: types.NewPerpetual(1, MarketId, sellerAcc.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(1000_000_000), math.LegacyZeroDec()), // Artificially high old margin
			trade:        types.NewTrade(MarketId, math.LegacyNewDec(5), math.LegacyNewDec(10), buyerAcc, sellerAcc),                                                        // Sell 5 @ 10 -> new avg ep = 70, new qty = -15
			isBuyer:      false,
			// Expecting error return
			expectedErrMsg: "newRequiredInitialMargin (105000000) must be greater than oldPerpetual.Margin (1000000000) for seller when position is increased from negative to more negative", // Match fmt.Errorf string
		},
	}

	for _, tc := range testCases {
		// Use t.Run for better test isolation and output if running multiple top-level tests
		// Or keep suite.Run if preferred structure
		suite.Run(tc.name, func() {
			// Reset balances to initial state for each run
			suite.SetAccountBalance(buyerAcc.GetTradingAccountAddress(), sdk.NewCoins(initialBalance))
			suite.SetAccountBalance(sellerAcc.GetTradingAccountAddress(), sdk.NewCoins(initialBalance))
			suite.SetAccountBalance(marketAccAddr, sdk.NewCoins(initialMarketBalanceCoin)) // Reset market balance too

			var primarySubAccount types.SubAccount
			var counterpartySubAccount types.SubAccount
			if tc.isBuyer {
				primarySubAccount = tc.trade.BuyerSubAccount
				counterpartySubAccount = tc.trade.SellerSubAccount
			} else {
				primarySubAccount = tc.trade.SellerSubAccount
				counterpartySubAccount = tc.trade.BuyerSubAccount
			}
			initialSubAccBalance := suite.GetAccountBalance(primarySubAccount.GetTradingAccountAddress(), QuoteDenom)
			initialMarketBalance := suite.GetAccountBalance(marketAccAddr, QuoteDenom)
			initialCounterpartyBalance := suite.GetAccountBalance(counterpartySubAccount.GetTradingAccountAddress(), QuoteDenom)

			// --- Execute ---
			perpetualToPass := tc.oldPerpetual // Pass copy
			updatedPerpetual, err := suite.app.ClobKeeper.SettleMarginAndRPnL(suite.ctx, market, perpetualToPass, false, tc.trade, tc.isBuyer)

			// --- Assert ---
			// Check if either an error return OR a panic was expected based on tc.expectedErrMsg existence or tc.expectPanic flag
			// Since we removed expectPanic flag again, just check based on expectedErrMsg
			if tc.expectedErrMsg != "" { // This now covers all expected failure modes (error returns)
				suite.Require().Error(err, "Expected an error but got nil for case: %s", tc.name)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg, "Error message mismatch for case: %s", tc.name)

				// Verify balances did not change on error
				finalSubAccBalance := suite.GetAccountBalance(primarySubAccount.GetTradingAccountAddress(), QuoteDenom)
				finalMarketBalance := suite.GetAccountBalance(marketAccAddr, QuoteDenom)
				finalCounterpartyBalance := suite.GetAccountBalance(counterpartySubAccount.GetTradingAccountAddress(), QuoteDenom)
				suite.Require().True(initialSubAccBalance.Equal(finalSubAccBalance), "SubAccount balance changed despite error. Initial: %s, Final: %s", initialSubAccBalance, finalSubAccBalance)
				suite.Require().True(initialMarketBalance.Equal(finalMarketBalance), "Market balance changed despite error. Initial: %s, Final: %s", initialMarketBalance, finalMarketBalance)
				suite.Require().True(initialCounterpartyBalance.Equal(finalCounterpartyBalance), "Counterparty balance changed despite error. Initial: %s, Final: %s", initialCounterpartyBalance, finalCounterpartyBalance)

			} else { // Expecting success
				suite.Require().NoError(err, "Expected no error but got: %v for case: %s", err, tc.name)

				// Check Perpetual State
				suite.Require().True(tc.expectedUpdatedQuantity.Equal(updatedPerpetual.Quantity), "Quantity mismatch. Expected %s, Got %s", tc.expectedUpdatedQuantity, updatedPerpetual.Quantity)
				if !updatedPerpetual.Quantity.IsZero() {
					suite.Require().True(tc.expectedUpdatedEntryPrice.Equal(updatedPerpetual.EntryPrice), "EntryPrice mismatch. Expected %s, Got %s", tc.expectedUpdatedEntryPrice, updatedPerpetual.EntryPrice)
					suite.Require().True(tc.expectedUpdatedMargin.Equal(updatedPerpetual.MarginAmount), "Margin mismatch. Expected %s, Got %s", tc.expectedUpdatedMargin, updatedPerpetual.MarginAmount)
				}

				// Check Balances
				finalSubAccBalance := suite.GetAccountBalance(primarySubAccount.GetTradingAccountAddress(), QuoteDenom)
				finalMarketBalance := suite.GetAccountBalance(marketAccAddr, QuoteDenom)
				finalCounterpartyBalance := suite.GetAccountBalance(counterpartySubAccount.GetTradingAccountAddress(), QuoteDenom)

				expectedFinalSubAccBalance := initialSubAccBalance.Add(tc.expectedSubAccBalanceChange)
				expectedFinalMarketBalance := initialMarketBalance.Add(tc.expectedMarketBalanceChange)

				suite.Require().True(expectedFinalSubAccBalance.Equal(finalSubAccBalance), "SubAccount balance mismatch. Initial %s, Change %s, Expected %s, Got %s", initialSubAccBalance, tc.expectedSubAccBalanceChange, expectedFinalSubAccBalance, finalSubAccBalance)
				suite.Require().True(expectedFinalMarketBalance.Equal(finalMarketBalance), "Market balance mismatch. Initial %s, Change %s, Expected %s, Got %s", initialMarketBalance, tc.expectedMarketBalanceChange, expectedFinalMarketBalance, finalMarketBalance)
				suite.Require().True(initialCounterpartyBalance.Equal(finalCounterpartyBalance), "Counterparty balance changed unexpectedly. Initial %s, Final: %s", initialCounterpartyBalance, finalCounterpartyBalance)
			}
		})
	}
}
