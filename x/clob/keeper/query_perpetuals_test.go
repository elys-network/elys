package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/clob/types"
)

func (suite *KeeperTestSuite) TestAllPerpetuals() {
	var market types.PerpetualMarket
	var liquidatingSubAccount types.SubAccount

	baseMmr := math.LegacyMustNewDecFromStr("0.05")

	setupTest := func() {
		market, liquidatingSubAccount, _, _ = suite.SetupExchangeTest()

		market.MaintenanceMarginRatio = baseMmr
		market.TwapPricesWindow = 3600
		suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)

		suite.FundAccount(liquidatingSubAccount.GetOwnerAccAddress(), sdk.NewCoins(sdk.NewInt64Coin(QuoteDenom, 1_000_000)))

		suite.SetTwapRecordDirectly(types.TwapPrice{
			MarketId:          MarketId,
			Block:             uint64(suite.ctx.BlockHeight() - 1),
			AverageTradePrice: math.LegacyNewDec(95),
			TotalVolume:       math.LegacyNewDec(1000),
			CumulativePrice:   math.LegacyNewDec(95000),
			Timestamp:         uint64(suite.ctx.BlockTime().Unix() - 60),
		})
	}

	testCases := []struct {
		name                      string
		perpetualsSetup           func() []types.Perpetual
		queryReq                  types.AllPerpetualsRequest
		expectedLiquidationPrices []math.LegacyDec
		expectedErrSubstring      string
	}{
		{
			name: "Success: Multiple perpetuals with valid liquidation prices",
			perpetualsSetup: func() []types.Perpetual {
				p1 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				p2 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				suite.SetPerpetualStateWithEntryFR(p1, false)
				suite.SetPerpetualStateWithEntryFR(p2, false)
				return []types.Perpetual{p1, p2}
			},
			queryReq: types.AllPerpetualsRequest{
				MarketId: 0,
				Pagination: &query.PageRequest{
					Limit: 2,
				},
			},
			expectedLiquidationPrices: []math.LegacyDec{math.LegacyMustNewDecFromStr("94.736842105263157895"), math.LegacyMustNewDecFromStr("104.761904761904761905")},
			expectedErrSubstring:      "",
		},
		{
			name: "Pagination: Limit 1, Offset 1",
			perpetualsSetup: func() []types.Perpetual {
				p1 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				p2 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				suite.SetPerpetualStateWithEntryFR(p1, false)
				suite.SetPerpetualStateWithEntryFR(p2, false)
				return []types.Perpetual{p1, p2}
			},
			queryReq: types.AllPerpetualsRequest{
				MarketId: 0,
				Pagination: &query.PageRequest{
					Limit:  1,
					Offset: 1,
				},
			},
			expectedLiquidationPrices: []math.LegacyDec{math.LegacyMustNewDecFromStr("104.761904761904761905")}, // The second perpetual's liquidation price
			expectedErrSubstring:      "",
		},
		{
			name: "Pagination: Limit 1, Offset 1, from specific market",
			perpetualsSetup: func() []types.Perpetual {
				p1 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				p2 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(-10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				suite.SetPerpetualStateWithEntryFR(p1, false)
				suite.SetPerpetualStateWithEntryFR(p2, false)
				return []types.Perpetual{p1, p2}
			},
			queryReq: types.AllPerpetualsRequest{
				MarketId: 1,
				Pagination: &query.PageRequest{
					Limit:  1,
					Offset: 1,
				},
			},
			expectedLiquidationPrices: []math.LegacyDec{math.LegacyMustNewDecFromStr("104.761904761904761905")}, // The second perpetual's liquidation price
			expectedErrSubstring:      "",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			setupTest()
			_ = tc.perpetualsSetup()

			res, err := suite.app.ClobKeeper.AllPerpetuals(suite.ctx, &tc.queryReq)

			if tc.expectedErrSubstring != "" {
				suite.Require().Error(err, "Expected an error but got nil")
				suite.Require().Contains(err.Error(), tc.expectedErrSubstring, "Error message mismatch")
			} else {
				suite.Require().NoError(err, "Expected no error but got: %v", err)
				suite.Require().Equal(len(tc.expectedLiquidationPrices), len(res.PerpetualInfos), "Mismatch in number of perpetuals returned")
				for i, expectedPrice := range tc.expectedLiquidationPrices {
					suite.Require().True(expectedPrice.Equal(res.PerpetualInfos[i].LiquidationPrice), "Mismatch in liquidation price for perpetual %d", i)
				}
			}
		})
	}
}
