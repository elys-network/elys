package keeper_test

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/x/clob/types"
)

func newTestPerpetualForLiqPrice(id uint64, owner string, qty math.LegacyDec, ep math.LegacyDec, marginAmt math.Int) types.Perpetual {
	return types.Perpetual{
		Id:           id,
		MarketId:     MarketId,
		Owner:        owner,
		Quantity:     qty,
		EntryPrice:   ep,
		MarginAmount: marginAmt,
		// EntryFundingRate is not used by GetLiquidationPrice
	}
}

func (suite *KeeperTestSuite) TestGetLiquidationPrice() {
	var subAccount types.SubAccount
	var market types.PerpetualMarket
	var perpetual types.Perpetual

	// Setup runs before each sub-test (t.Run)
	setupSubTest := func() {
		// Use a fresh context for each sub-test to avoid mock call count issues
		// and ensure clean state if any other keepers were involved.
		// suite.ctx = suite.app.BaseApp.NewContext(false, tmproto.Header{Height: suite.ctx.BlockHeight() + 1}) // Example

		subAccount = types.SubAccount{
			Owner:    "test_owner",
			MarketId: MarketId,
			// Balances not directly used by GetLiquidationPrice for isolated margin
		}
		market = types.PerpetualMarket{
			Id:                     MarketId,
			QuoteDenom:             QuoteDenom,
			MaintenanceMarginRatio: math.LegacyMustNewDecFromStr("0.05"), // 5% MMR
			// Other market fields
		}
		// Default perpetual (long)
		perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000)) // 10 units @ EP 100, Margin 100
	}

	testCases := []struct {
		name                 string
		setupSpecific        func() // To modify default subAccount, market, perpetual, or mock Oracle
		expectedLiqPriceStr  string // Expected price as string for precise LegacyDec comparison, or "" if error
		expectedErrSubstring string // Substring of expected error, or "" for no error
	}{
		// --- Success Cases (Isolated Margin) ---
		{
			name: "Success: Long position, standard parameters",
			setupSpecific: func() {
				// Uses default setup: Long Qty=10, EP=100, Margin=100, MMR=0.05, QuotePrice=1
				// IMV = 100 * 1 = 100
				// Num = 100 - (10 * 100) = -900
				// Den = (0.05 * 10) - 10 = 0.5 - 10 = -9.5
				// LiqPx = -900 / -9.5 = 94.736842105263157895
			},
			expectedLiqPriceStr:  "94.736842105263157895",
			expectedErrSubstring: "",
		},
		{
			name: "Success: Short position, standard parameters",
			setupSpecific: func() {
				perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000)) // Short Qty=-10
				// IMV = 100 * 1 = 100
				// Num = 100 - (-10 * 100) = 100 - (-1000) = 1100
				// Den = (0.05 * abs(-10)) - (-10) = (0.05 * 10) + 10 = 0.5 + 10 = 10.5
				// LiqPx = 1100 / 10.5 = 104.761904761904761905
			},
			expectedLiqPriceStr:  "104.761904761904761905",
			expectedErrSubstring: "",
		},
		{
			name: "Success: Long position, higher margin added",
			setupSpecific: func() {
				perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(200_000_000)) // Margin 200
				// IMV = 200 * 1 = 200
				// Num = 200 - (10 * 100) = -800
				// Den = (0.05 * 10) - 10 = -9.5
				// LiqPx = -800 / -9.5 = 84.210526315789473684
			},
			expectedLiqPriceStr:  "84.210526315789473684",
			expectedErrSubstring: "",
		},
		{
			name: "Success: Long position, zero margin (bankruptcy price)",
			setupSpecific: func() {
				perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.ZeroInt()) // Margin 0
				// IMV = 0 * 1 = 0
				// Num = 0 - (10 * 100) = -1000
				// Den = (0.05 * 10) - 10 = -9.5
				// LiqPx = -1000 / -9.5 = 105.263157894736842105 (Bankruptcy Price for Long if starting with 0 margin)
			},
			expectedLiqPriceStr:  "105.263157894736842105",
			expectedErrSubstring: "",
		},
		{
			name: "Success: Quote Denom Price is not 1",
			setupSpecific: func() {
				// Long Qty=10, EP=100, MarginAmount=100 tokens, MMR=0.05
				// QuoteDenomPrice = 0.5 (e.g., token is worth $0.5)
				// IMV = 100 * 0.5 = 50
				// Num = 50 - (10 * 100) = -950
				// Den = (0.05 * 10) - 10 = -9.5
				// LiqPx = -950 / -9.5 = 100
				suite.SetPrice([]string{"USDC"}, []math.LegacyDec{math.LegacyMustNewDecFromStr("0.5")})
			},
			expectedLiqPriceStr:  "100.000000000000000000",
			expectedErrSubstring: "",
		},

		// --- Error Cases (Isolated Margin) ---
		{
			name: "Error: Denominator zero (MMR = 1 for Long)",
			setupSpecific: func() {
				market.MaintenanceMarginRatio = math.LegacyOneDec()                                                                                   // MMR = 1
				perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000)) // Long
			},
			expectedErrSubstring: "division by zero (MMR = 1)",
		},
		{
			name: "Error: Denominator zero (MMR = -1 for Short - though MMR should be positive)",
			setupSpecific: func() {
				// This case tests the formula, though MMR=-1 is not a valid market param
				market.MaintenanceMarginRatio = math.LegacyOneDec().Neg()                                                                              // MMR = -1
				perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000)) // Short
			},
			// Den = (MMR * AbsQty) - Qty = (-1 * 10) - (-10) = -10 + 10 = 0
			expectedErrSubstring: "division by zero", // More generic as MMR=-1 is not the specific "MMR = 1" case
		},
		{
			name: "Error: Zero quantity (denominator becomes zero)",
			setupSpecific: func() {
				perpetual = newTestPerpetualForLiqPrice(1, subAccount.Owner, math.LegacyZeroDec(), math.LegacyNewDec(100), math.NewInt(100_000_000)) // Zero Qty
			},
			// Den = (MMR * 0) - 0 = 0
			expectedErrSubstring: "division by zero",
		},
		// Should be tested in end
		{
			name: "Error: GetDenomPrice fails",
			setupSpecific: func() {
				suite.app.OracleKeeper.RemovePrice(suite.ctx, "USDC", "test", uint64(suite.ctx.BlockTime().Unix()))
			},
			expectedErrSubstring: "denom (uusdc) price not found",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			setupSubTest() // Reset defaults for subAccount, market, perpetual

			// Apply test-case specific modifications
			if tc.setupSpecific != nil {
				tc.setupSpecific()
			}

			// Define the call
			call := func() (math.LegacyDec, error) {
				// Pass copies to ensure original tc.perpetual/market are not modified if they are pointers
				pCopy := perpetual
				mCopy := market
				saCopy := subAccount
				return suite.app.ClobKeeper.GetLiquidationPrice(suite.ctx, pCopy, mCopy, saCopy)
			}

			liqPrice, err := call()

			if tc.expectedErrSubstring != "" {
				suite.Require().Error(err, "Expected an error but got nil")
				suite.Require().Contains(err.Error(), tc.expectedErrSubstring, "Error message mismatch")
			} else {
				suite.Require().NoError(err, "Expected no error but got: %v", err)
				suite.Require().NotNil(liqPrice, "Liquidation price should not be nil on success")
				if !liqPrice.IsNil() { // Check before calling methods on it
					expectedPrice := math.LegacyMustNewDecFromStr(tc.expectedLiqPriceStr)
					// Compare strings for exact LegacyDec match due to potential precision nuances
					suite.Require().Equal(expectedPrice.String(), liqPrice.String(), "Liquidation price mismatch. Expected %s, Got %s", expectedPrice.String(), liqPrice.String())
				}
			}
		})
	}
}
