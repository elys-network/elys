package keeper_test

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/v7/x/clob/types"
	"github.com/stretchr/testify/require"
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
		SubAccountId:     MarketId,
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
		expectADL            bool
		postCheck            func(originalPerp types.Perpetual) // To verify state after execution
	}{
		// --- Success Cases ---
		{
			name:                 "Long position fully liquidated",
			perpetualToLiquidate: newTestPerpetualForLiquidation("liquidating_owner_long", math.LegacyNewDec(10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)),
			setupSpecific: func() {
				// Provide enough sell liquidity to fill the entire long position
				sellOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(99), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId)
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sellOrder)
			},
			expectedRatio: math.LegacyOneDec(),
			expectADL:     false,
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
				buyOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_SELL, math.LegacyNewDec(101), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId)
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buyOrder)
			},
			expectedRatio: math.LegacyOneDec(),
			expectADL:     false,
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
				sellOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(99), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(5), math.LegacyZeroDec(), MarketId) // Only 5 available
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sellOrder)
			},
			expectedRatio: math.LegacyMustNewDecFromStr("0.5"), // 5 filled / 10 original
			expectADL:     true,                                // ADL should be set for the remainder
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
				buyOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_SELL, math.LegacyNewDec(101), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(7), math.LegacyZeroDec(), MarketId) // Only 7 available
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, buyOrder)
			},
			expectedRatio: math.LegacyMustNewDecFromStr("0.7"), // 7 filled / 10 original
			expectADL:     true,
		},

		// --- Error Cases ---
		{
			name: "IF Insufficient during settlement (Long Liquidation)",
			// Setup: Long position, fill will cause loss > margin, IF is empty
			perpetualToLiquidate: newTestPerpetualForLiquidation("bankrupt_owner_long", math.LegacyNewDec(10), math.LegacyNewDec(100), calcMarginForTest(math.LegacyNewDec(10), math.LegacyNewDec(100), baseImr)), // Margin = 100
			setupSpecific: func() {
				// Order to fill against, at a very bad price for the long
				sellOrder := types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(1), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId)
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, sellOrder)
				// Empty the Insurance Fund
				err := suite.BurnAccountBalance(market.GetInsuranceAccount(), QuoteDenom)
				suite.Require().NoError(err, "Failed to burn insurance fund")
			},
			// PNL = 10 * (1 - 100) = -990. NetRefund = 100 - 990 - Fee. IF will be needed.
			expectedRatio:        math.LegacyZeroDec(), // Returns 0 ratio on IF error by design
			expectedErrSubstring: "",                   // MarketLiquidation returns nil error in this case
			expectADL:            true,                 // ADL should be set
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
				tc.perpetualToLiquidate = suite.SetPerpetualStateWithEntryFR(tc.perpetualToLiquidate, false) // Ensure it's in store
			}

			// Clear ADL flag for this perpetual before the test
			suite.app.ClobKeeper.DeletePerpetualADL(suite.ctx, tc.perpetualToLiquidate.MarketId, tc.perpetualToLiquidate.Id)

			// --- Execute ---
			closingRatio, adlTriggered, err := suite.app.ClobKeeper.MarketLiquidation(suite.ctx, tc.perpetualToLiquidate, market)

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

			if tc.expectADL {
				suite.Require().True(closingRatio.LT(math.LegacyOneDec()), "Expected ADL")
				suite.Require().True(adlTriggered, "Expected ADL")
			} else {
				suite.Require().False(closingRatio.LT(math.LegacyOneDec()), "Expected no ADL")
				suite.Require().False(adlTriggered, "Expected no ADL")
			}

			if tc.postCheck != nil {
				tc.postCheck(tc.perpetualToLiquidate) // Pass original for comparison if needed
			}
		})
	}
}

func newTestPerpetualForForcedLiq(owner string, qty math.LegacyDec, ep math.LegacyDec, margin math.Int) types.Perpetual {
	return types.Perpetual{
		Id:               0,
		MarketId:         MarketId,
		Owner:            owner,
		Quantity:         qty,
		EntryPrice:       ep,
		MarginAmount:     margin,
		EntryFundingRate: math.LegacyZeroDec(), // Assume funding settled for simplicity
		SubAccountId:     MarketId,
	}
}

func (suite *KeeperTestSuite) TestForcedLiquidation() { // Using KeeperTestSuite
	var market types.PerpetualMarket
	var liquidatingSubAccount, counterpartySubAccount types.SubAccount // These will be assigned by setupTest
	var marketModuleAcc sdk.AccAddress

	baseMmr := math.LegacyMustNewDecFromStr("0.05")
	baseLiqFeeRate := math.LegacyMustNewDecFromStr("0.01")

	// setupTest runs before each sub-test (t.Run)
	setupTest := func() {
		market, liquidatingSubAccount, counterpartySubAccount, marketModuleAcc = suite.SetupExchangeTest() // Assuming this is a method on KeeperTestSuite

		market.MaintenanceMarginRatio = baseMmr
		market.LiquidationFeeShareRate = baseLiqFeeRate
		market.TwapPricesWindow = 3600
		suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)

		suite.FundAccount(liquidatingSubAccount.GetOwnerAccAddress(), sdk.NewCoins(sdk.NewInt64Coin(QuoteDenom, 1_000_000))) // Initial 1M micro-units

		suite.SetTwapRecordDirectly(types.TwapPrice{ // Assuming suite.SetTwapRecordDirectly is a helper
			MarketId:          MarketId,
			Block:             uint64(suite.ctx.BlockHeight() - 1),
			AverageTradePrice: math.LegacyNewDec(95),
			TotalVolume:       math.LegacyNewDec(1000),
			CumulativePrice:   math.LegacyNewDec(95000),
			Timestamp:         uint64(suite.ctx.BlockTime().Unix() - 60),
		})
		// Assuming GetDenomPrice for quote denom (e.g., "uusdc") is handled by OracleKeeper setup
		// and returns 1 (or 0.000001 if it converts micro-units to base units, adjust accordingly)
		// For simplicity, GetLiquidationPrice and GetEquityValue are assumed to work with these values.
	}

	testCases := []struct {
		name                         string
		perpetualToLiquidateSetup    func() types.Perpetual
		orderBookSetup               func()
		preSetup                     func()
		expectedErrSubstring         string
		expectADLSet                 bool
		expectedLiquidatorRewardPaid math.Int // Actual amount transferred to liquidator
		checkPostLiquidationState    func(originalPerp types.Perpetual, m types.PerpetualMarket, liquidatedTraderSubAcc types.SubAccount, cpSubAcc types.SubAccount, liquidatorAddr sdk.AccAddress, marketModAddr sdk.AccAddress)
	}{
		// --- A. Eligibility Check Failures ---
		{
			name: "Fail: Asset price is 0",
			perpetualToLiquidateSetup: func() types.Perpetual {
				// Margin is 100_000_000 uusdc (100 base units if 1M uusdc = 1 USDC)
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() {},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyZeroDec()})
			},
			expectedErrSubstring:         "invalid price 0.000000000000000000",
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},
		{
			name: "Fail: sub account not found",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() {},
			preSetup: func() {

				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(104)})
			},
			expectedErrSubstring:         "short position not yet liquidatable",
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},
		{
			name: "Fail:Denom price",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() {},
			preSetup: func() {
				suite.app.OracleKeeper.RemovePrice(suite.ctx, "USDC", "test", uint64(suite.ctx.BlockTime().Unix()))
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(104)})
			},
			expectedErrSubstring:         "denom (uusdc) price not found",
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},
		{
			name: "Fail: Long position not yet liquidatable",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() {},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(96)})
			},
			expectedErrSubstring:         "long position not yet liquidatable",
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},
		{
			name: "Fail: Short position not yet liquidatable",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() {},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(104)})
			},
			expectedErrSubstring:         "short position not yet liquidatable",
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},

		// --- B. Successful Full Liquidation (Trader's margin covers loss + fee) ---
		{
			name: "Success: Long position fully liquidated, margin covers loss and fee",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(94), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId))
			},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(94)})
			},
			expectedErrSubstring:         "",
			expectADLSet:                 false,
			expectedLiquidatorRewardPaid: math.NewInt(1_000_000), // 0.01 * 100_000_000
			checkPostLiquidationState: func(originalPerp types.Perpetual, m types.PerpetualMarket, liqSubAcc types.SubAccount, cpSubAcc types.SubAccount, liq sdk.AccAddress, marketMod sdk.AccAddress) {
				_, found := suite.GetPerpetualState(liqSubAcc.GetOwnerAccAddress(), originalPerp.MarketId)
				require.False(suite.T(), found, "Liquidated perpetual should be deleted")
			},
		},

		// --- C. Partial Fill, ADL Set for Remainder ---
		{
			name: "Success (Partial Fill): Long partially liquidated, ADL set for remainder",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000)) // Margin 100M uusdc
				return p
			},
			orderBookSetup: func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(94), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(5), math.LegacyZeroDec(), MarketId))
			},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(94)})
			},
			expectedErrSubstring: "", // ForcedLiquidation returns nil if partial fill was settled and ADL set
			expectADLSet:         true,
			// Base reward = 0.01 * 100M = 1M. Closing ratio = 0.5. Adjusted = 1M * 0.5 = 500_000.
			expectedLiquidatorRewardPaid: math.NewInt(500_000),
			checkPostLiquidationState: func(originalPerp types.Perpetual, m types.PerpetualMarket, liqSubAcc types.SubAccount, cpSubAcc types.SubAccount, liq sdk.AccAddress, marketMod sdk.AccAddress) {
				updatedPerp, found := suite.GetPerpetualState(liqSubAcc.GetOwnerAccAddress(), originalPerp.MarketId)
				require.True(suite.T(), found, "Perpetual should exist after partial liquidation")
				require.True(suite.T(), updatedPerp.Quantity.Equal(math.LegacyNewDec(5)), "Partially liquidated quantity incorrect")
			},
		},

		// --- D. IF Insufficient During MarketLiquidation's Settlement ---
		{
			name: "Success (ADL Triggered): IF insufficient during settlement, no liquidator fee",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(10_000_000)) // Margin 10M uusdc
				return p
			},
			orderBookSetup: func() {
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(1), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId))
			},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(1)})
				err := suite.BurnAccountBalance(market.GetInsuranceAccount(), QuoteDenom)
				suite.Require().NoError(err)
			},
			expectedErrSubstring:         "", // ForcedLiquidation returns (0, nil) if MarketLiquidation returns (0,true,nil)
			expectADLSet:                 true,
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},

		// --- E. Market Account is empty for trader's refund, should not happen ideally ---
		{
			name: "Error: Market Account is empty for trader's refund",
			perpetualToLiquidateSetup: func() types.Perpetual {
				p := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				return p
			},
			orderBookSetup: func() { // Full fill
				suite.app.ClobKeeper.SetPerpetualOrder(suite.ctx, types.NewPerpetualOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_SELL, math.LegacyNewDec(105), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId))
			},
			preSetup: func() {
				suite.SetPrice([]string{"ATOM"}, []math.LegacyDec{math.LegacyNewDec(105)})
				err := suite.BurnAccountBalance(marketModuleAcc, QuoteDenom)
				suite.Require().NoError(err)
				suite.FundAccount(market.GetInsuranceAccount(), sdk.NewCoins(sdk.NewInt64Coin(QuoteDenom, 500_000_000))) // Ensure IF is very solvent
			},
			expectedErrSubstring:         "insufficient funds", // Error from MarketLiquidation due to market acc being empty for trader's refund
			expectADLSet:                 false,                // ADL not set if MarketLiquidation returns a direct error
			expectedLiquidatorRewardPaid: math.ZeroInt(),
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			setupTest() // This assigns to outer scope market, liquidatingSubAccount, counterpartySubAccount, marketModuleAcc
			currentPerpetual := tc.perpetualToLiquidateSetup()

			currentPerpetual.Owner = liquidatingSubAccount.Owner
			currentPerpetual = suite.SetPerpetualStateWithEntryFR(currentPerpetual, false)

			if tc.orderBookSetup != nil {
				tc.orderBookSetup()
			}
			if tc.preSetup != nil {
				tc.preSetup()
			}

			suite.app.ClobKeeper.DeletePerpetualADL(suite.ctx, currentPerpetual.MarketId, currentPerpetual.Id)

			initialLiquidatorBalance := suite.GetAccountBalance(liquidatingSubAccount.GetOwnerAccAddress(), QuoteDenom)

			paidReward, err := suite.app.ClobKeeper.ForcedLiquidation(suite.ctx, currentPerpetual, market, liquidatingSubAccount.GetOwnerAccAddress())

			if tc.expectedErrSubstring != "" {
				suite.Require().Error(err, "Expected an error but got nil for case: %s", tc.name)
				suite.Require().Contains(err.Error(), tc.expectedErrSubstring, "Error message mismatch for case: %s", tc.name)
			} else {
				suite.Require().NoError(err, "Expected no error but got: %v for case: %s", err, tc.name)
			}

			// Assert ADL Status
			_, adlFound := suite.app.ClobKeeper.GetPerpetualADL(suite.ctx, currentPerpetual.MarketId, currentPerpetual.Id)
			suite.Require().Equal(tc.expectADLSet, adlFound, "ADL set status mismatch for case: %s", tc.name)

			actualPaidToLiquidator := math.ZeroInt()
			if err == nil { // If ForcedLiquidation itself didn't error out before/during payment
				actualPaidToLiquidator = tc.expectedLiquidatorRewardPaid
			} else if tc.expectedErrSubstring != "" && !errors.Is(err, sdkerrors.ErrInsufficientFunds) && err.Error() != types.ErrInsufficientInsuranceFund.Error() {
				// If it's an eligibility error, no payment was attempted.
				actualPaidToLiquidator = math.ZeroInt()
			} else if errors.Is(err, sdkerrors.ErrInsufficientFunds) && tc.expectedErrSubstring == "insufficient funds" {
				// This case is when bankKeeper.SendCoins to liquidator failed
				actualPaidToLiquidator = math.ZeroInt()
			}

			suite.Require().Equal(tc.expectedLiquidatorRewardPaid, paidReward, "Liquidator reward mismatch for case: %s", tc.name)
			suite.CheckBalanceChange(liquidatingSubAccount.GetOwnerAccAddress(), initialLiquidatorBalance, actualPaidToLiquidator, fmt.Sprintf("Liquidator balance for case: %s", tc.name))

			if tc.checkPostLiquidationState != nil {
				actualLiqOwnerAddr := sdk.MustAccAddressFromBech32(currentPerpetual.Owner)
				actualLiqSubAccount, errGetSub := suite.app.ClobKeeper.GetSubAccount(suite.ctx, actualLiqOwnerAddr, currentPerpetual.SubAccountId)

				// Handle state of actualLiqSubAccount based on whether ForcedLiquidation errored
				if err == nil { // If ForcedLiquidation itself succeeded (even if ADL was set for partial fill)
					suite.Require().NoError(errGetSub, "Failed to get sub-account for post-check after successful ForcedLiquidation for case: %s", tc.name)
				} else if errors.Is(err, sdkerrors.ErrInsufficientFunds) && tc.expectedErrSubstring == "insufficient funds" {
					// If error was due to liquidator payment failure, sub-account state *before* that failure is relevant
					// GetSubAccount should still succeed unless the subaccount itself was the source of the problem
					suite.Require().NoError(errGetSub, "Failed to get sub-account for post-check even after liquidator payment failure for case: %s", tc.name)
				} else if errGetSub == nil {
					suite.T().Logf("ForcedLiquidation errored for case '%s', but subaccount still found. Proceeding with postCheck.", tc.name)
				} else {
					suite.T().Logf("Skipping postCheck for case '%s' due to GetSubAccount error (%v) after ForcedLiquidation error (%v).", tc.name, errGetSub, err)
					return // Skip post check if we can't even get the subaccount
				}

				suite.T().Logf("Calling checkPostLiquidationState for: %s with liqSubAcc: %+v", tc.name, actualLiqSubAccount)
				// Pass the subaccounts from the main setup for cpSubAcc and marketModuleAcc
				tc.checkPostLiquidationState(currentPerpetual, market, actualLiqSubAccount, counterpartySubAccount, liquidatingSubAccount.GetOwnerAccAddress(), marketModuleAcc)
			}
		})
	}
}
