package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

// Helper to create a basic perpetual for tests
func newTestPerpetualForLiquidation(owner string, qty math.LegacyDec, ep math.LegacyDec, margin math.Int) types.Perpetual {
	return types.Perpetual{
		Id:               0,
		MarketId:         MarketId,
		Owner:            owner,
		Quantity:         qty,
		EntryPrice:       ep,
		MarginAmount:     margin,
		EntryFundingRate: math.LegacyZeroDec(), // Or set from current market rate if relevant
	}
}

// Local calcMargin for test case definition clarity
func calcMarginForTest(qty, price, imr math.LegacyDec) math.Int {
	if price.IsNil() || !price.IsPositive() || qty.IsZero() {
		return math.ZeroInt()
	}
	return qty.Abs().Mul(price).Mul(imr).MulInt64(1_000_000).RoundInt()
}

func (suite *KeeperTestSuite) TestMarketLiquidation() {
	var market types.PerpetualMarket
	var liquidatingTraderSubAccount, counterpartySubAccount types.SubAccount

	baseImr := math.LegacyMustNewDecFromStr("0.1") // 10% IMR

	// Default setup for each subtest run
	defaultSetup := func() {
		market, liquidatingTraderSubAccount, counterpartySubAccount, _ = suite.SetupExchangeTest()
		baseImr = market.InitialMarginRatio
		// Ensure TWAP returns a non-zero price for most tests
		suite.SetTwapPriceDirectly(market.Id, math.LegacyNewDec(100)) // Assume helper sets this mock/state
		// Ensure initial funding rate is set for perpetual state consistency
		suite.app.ClobKeeper.SetFundingRate(suite.ctx, types.FundingRate{MarketId: MarketId, Rate: math.LegacyZeroDec(), Block: uint64(suite.ctx.BlockHeight())})
	}

	testCases := []struct {
		name                 string
		perpetualToLiquidate types.Perpetual
		setupSpecific        func() // To set up liquidity, IF balance, etc.
		expectedRatio        math.LegacyDec
		expectedErrSubstring string // Substring of expected error, or "" for no error
		expectADLSet         bool
		postCheck            func(originalPerp types.Perpetual) // To verify state after execution
	}{
		// --- Success Cases ---
		{
			name:                 "Long position fully liquidated",
			perpetualToLiquidate: newTestPerpetualForLiquidation("liquidating_owner_long", math.LegacyNewDec(10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)),
			setupSpecific: func() {
				// Provide enough sell liquidity to fill the entire long position
				sellOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(99), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec())
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sellOrder)
			},
			expectedRatio: math.LegacyOneDec(),
			expectADLSet:  false,
			postCheck: func(originalPerp types.Perpetual) {
				_, found := suite.GetPerpetualState(sdk.MustAccAddressFromBech32(originalPerp.Owner), originalPerp.MarketId)
				suite.Require().False(found, "Perpetual should be deleted after full liquidation")
			},
		},
		{
			name:                 "Short position fully liquidated",
			perpetualToLiquidate: newTestPerpetualForLiquidation("liquidating_owner_short", math.LegacyNewDec(-10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)),
			setupSpecific: func() {
				// Provide enough buy liquidity
				buyOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_SELL, math.LegacyNewDec(101), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec())
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buyOrder)
			},
			expectedRatio: math.LegacyOneDec(),
			expectADLSet:  false,
			postCheck: func(originalPerp types.Perpetual) {
				_, found := suite.GetPerpetualState(sdk.MustAccAddressFromBech32(originalPerp.Owner), originalPerp.MarketId)
				suite.Require().False(found, "Perpetual should be deleted after full liquidation")
			},
		},
		{
			name:                 "Long position partially liquidated due to insufficient sell liquidity",
			perpetualToLiquidate: newTestPerpetualForLiquidation("liquidating_owner_long_partial", math.LegacyNewDec(10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)),
			setupSpecific: func() {
				// Provide less sell liquidity than needed
				sellOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(99), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(5), math.LegacyZeroDec()) // Only 5 available
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sellOrder)
			},
			expectedRatio: math.LegacyMustNewDecFromStr("0.5"), // 5 filled / 10 original
			expectADLSet:  true,                                // ADL should be set for the remainder
			postCheck: func(originalPerp types.Perpetual) {
				updatedPerp, found := suite.GetPerpetualState(sdk.MustAccAddressFromBech32(originalPerp.Owner), originalPerp.MarketId)
				suite.Require().True(found, "Perpetual should still exist after partial liquidation")
				// Original Qty +10, 5 filled. Remaining Qty +5.
				expectedRemainingQty := math.LegacyNewDec(5)
				suite.Require().True(updatedPerp.Quantity.Equal(expectedRemainingQty), "Perpetual quantity not correctly reduced. Expected %s, Got %s", expectedRemainingQty, updatedPerp.Quantity)
			},
		},
		{
			name:                 "Short position partially liquidated due to insufficient buy liquidity",
			perpetualToLiquidate: newTestPerpetualForLiquidation("liquidating_owner_short_partial", math.LegacyNewDec(-10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)),
			setupSpecific: func() {
				buyOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_SELL, math.LegacyNewDec(101), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(7), math.LegacyZeroDec()) // Only 7 available
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buyOrder)
			},
			expectedRatio: math.LegacyMustNewDecFromStr("0.7"), // 7 filled / 10 original
			expectADLSet:  true,
		},

		// --- Error Cases ---
		{
			name: "IF Insufficient during settlement (Long Liquidation)",
			// Setup: Long position, fill will cause loss > margin, IF is empty
			perpetualToLiquidate: newTestPerpetualForLiquidation("bankrupt_owner_long", math.LegacyNewDec(10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)), // Margin = 100
			setupSpecific: func() {
				// Order to fill against, at a very bad price for the long
				sellOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(1), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec())
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sellOrder)
				// Empty the Insurance Fund
				err := suite.BurnAccountBalance(market.GetInsuranceAccount(), QuoteDenom)
				suite.Require().NoError(err, "Failed to burn insurance fund")
			},
			// PNL = 10 * (1 - 100) = -990. NetRefund = 100 - 990 - Fee. IF will be needed.
			expectedRatio:        math.LegacyZeroDec(), // Returns 0 ratio on IF error by design
			expectedErrSubstring: "",                   // MarketLiquidation returns nil error in this case
			expectADLSet:         true,                 // ADL should be set
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			defaultSetup()
			if tc.setupSpecific != nil {
				tc.setupSpecific()
			}

			// Ensure perpetual's owner matches one of the subaccounts for GetSubAccount to work if called for this perp
			// This is more relevant for ForcedLiquidation, MarketLiquidation uses msg.Creator
			tc.perpetualToLiquidate.Owner = liquidatingTraderSubAccount.Owner
			// For this test, the actual liquidatingTraderSubAccount doesn't need its position state set,
			// as that's handled by `Exchange` called internally by `ExecuteMarket...Order`.
			// We need to save the `perpetualToLiquidate` to the store if it's supposed to be fetched by internal logic.
			// However, `MarketLiquidation` takes `perpetual` as an argument.
			// Let's assume it's okay if it's not in store before, `Exchange` handles it.
			// But `GetPerpetual` *within* MarketLiquidation (after partial fill) needs it. So it should be set initially.
			if !tc.perpetualToLiquidate.Quantity.IsZero() { // Only set if it's a defined position
				tc.perpetualToLiquidate = suite.SetPerpetualStateWithEntryFR(tc.perpetualToLiquidate) // Ensure it's in store
			}

			// Clear ADL flag for this perpetual before the test
			suite.app.ClobKeeper.DeletePerpetualADL(suite.ctx, tc.perpetualToLiquidate.MarketId, tc.perpetualToLiquidate.Id)

			// --- Execute ---
			closingRatio, err := suite.app.ClobKeeper.MarketLiquidation(suite.ctx, tc.perpetualToLiquidate, market)

			// --- Assert ---
			if tc.expectedErrSubstring != "" {
				suite.Require().Error(err, "Expected an error but got nil")
				suite.Require().Contains(err.Error(), tc.expectedErrSubstring, "Error message mismatch")
			} else {
				suite.Require().NoError(err, "Expected no error but got: %v", err)
				suite.Require().NotNil(closingRatio, "ClosingRatio should not be nil on success")
				if !closingRatio.IsNil() { // Avoid panic on nil comparison if NoError path still returns nil ratio
					suite.Require().True(tc.expectedRatio.Equal(closingRatio), "ClosingRatio mismatch. Expected %s, Got %s", tc.expectedRatio, closingRatio)
				}
			}

			// Check ADL status
			adlRecord, adlFound := suite.app.ClobKeeper.GetPerpetualADL(suite.ctx, tc.perpetualToLiquidate.MarketId, tc.perpetualToLiquidate.Id)
			if tc.expectADLSet {
				suite.Require().True(adlFound, "Expected ADL record to be set, but not found")
				if adlFound { // Avoid nil dereference if not found
					suite.Require().Equal(tc.perpetualToLiquidate.Id, adlRecord.Id, "ADL record ID mismatch")
					suite.Require().Equal(tc.perpetualToLiquidate.MarketId, adlRecord.MarketId, "ADL record MarketId mismatch")
				}
			} else {
				suite.Require().False(adlFound, "Expected no ADL record, but one was found")
			}

			if tc.postCheck != nil {
				tc.postCheck(tc.perpetualToLiquidate) // Pass original for comparison if needed
			}
		})
	}
}

// --- Helper stubs assumed to exist in suite ---
/*
func (suite *ExchangeTestSuite) SetTwapPriceDirectly(marketId uint64, price math.LegacyDec) {
    // Mocking or direct store manipulation. For this test, we'd assume
    // GetCurrentTwapPrice is called *within* Exchange -> SettleMarginAndRPnL -> OnPositionClose,
    // so its behavior is part of what ExecuteMarketBuy/SellOrder implicitly tests.
    // The TWAP price set in setupDefault directly affects how OnPositionClose calculates PNL.
}
func (suite *ExchangeTestSuite) DeletePerpetualADL(ctx sdk.Context, marketId uint64, perpetualId uint64) {
    // ...
}
func (suite *ExchangeTestSuite) GetPerpetualADL(ctx sdk.Context, marketId uint64, perpetualId uint64) (types.PerpetualADL, bool) {
    // ...
}
// ... other required helpers: SetupExchangeTest, CreateMarketWithIMR, SetPerpetualState, GetPerpetualState, BurnAccountBalance ...
*/
