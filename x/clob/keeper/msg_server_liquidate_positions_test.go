package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (suite *KeeperTestSuite) TestLiquidatePositions() {
	var market types.PerpetualMarket
	var liquidatingSubAccount, counterpartySubAccount types.SubAccount

	baseMmr := math.LegacyMustNewDecFromStr("0.05")
	baseLiqFeeRate := math.LegacyMustNewDecFromStr("0.01")

	setupTest := func() {
		market, liquidatingSubAccount, counterpartySubAccount, _ = suite.SetupExchangeTest()

		market.MaintenanceMarginRatio = baseMmr
		market.LiquidationFeeShareRate = baseLiqFeeRate
		market.TwapPricesWindow = 3600
		suite.app.ClobKeeper.SetPerpetualMarket(suite.ctx, market)

		suite.FundAccount(liquidatingSubAccount.GetOwnerAccAddress(), sdk.NewCoins(sdk.NewInt64Coin(QuoteDenom, 1_000_000)))

		suite.SetTwapRecordDirectly(types.TwapPrice{
			MarketId:          MarketId,
			Block:             uint64(suite.ctx.BlockHeight() - 1),
			AverageTradePrice: math.LegacyNewDec(94),
			TotalVolume:       math.LegacyNewDec(1000),
			CumulativePrice:   math.LegacyNewDec(95000),
			Timestamp:         uint64(suite.ctx.BlockTime().Unix() - 60),
		})
	}

	testCases := []struct {
		name                 string
		positionsSetup       func() []types.LiquidatePosition
		orderBookSetup       func()
		expectedReward       sdk.Coins
		expectedErrSubstring string
	}{
		{
			name: "Success: Liquidate multiple positions",
			positionsSetup: func() []types.LiquidatePosition {
				p1 := newTestPerpetualForForcedLiq(liquidatingSubAccount.Owner, math.LegacyNewDec(10), math.LegacyNewDec(100), math.NewInt(100_000_000))
				suite.SetPerpetualStateWithEntryFR(p1, false)
				return []types.LiquidatePosition{
					{MarketId: 1, PerpetualId: 1},
				}
			},
			orderBookSetup: func() {
				suite.app.ClobKeeper.SetOrder(suite.ctx, types.NewOrder(MarketId, types.OrderType_ORDER_TYPE_LIMIT_BUY, math.LegacyNewDec(94), uint64(suite.ctx.BlockHeight()), counterpartySubAccount.GetOwnerAccAddress(), math.LegacyNewDec(10), math.LegacyZeroDec(), MarketId))
			},
			expectedReward:       sdk.NewCoins(sdk.NewCoin(QuoteDenom, math.NewInt(1_000_000))), // 1M reward for each position
			expectedErrSubstring: "",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			setupTest()
			positions := tc.positionsSetup()
			if tc.orderBookSetup != nil {
				tc.orderBookSetup()
			}

			msg := &types.MsgLiquidatePositions{
				Liquidator: liquidatingSubAccount.GetOwnerAccAddress().String(),
				Positions:  positions,
			}

			_, err := suite.app.ClobKeeper.LiquidatePositions(suite.ctx, msg)

			if tc.expectedErrSubstring != "" {
				suite.Require().Error(err, "Expected an error but got nil")
				suite.Require().Contains(err.Error(), tc.expectedErrSubstring, "Error message mismatch")
			} else {
				suite.Require().NoError(err, "Expected no error but got: %v", err)
			}
		})
	}
}
