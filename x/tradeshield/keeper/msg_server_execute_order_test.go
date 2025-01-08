package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	keeper "github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (suite *TradeshieldKeeperTestSuite) TestMsgServerExecuteOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                  string
		expectErrMsg          string
		prerequisiteFunction  func() *types.MsgExecuteOrders
		postrequisiteFunction func()
	}{
		{
			"Spot Order not found",
			"spot order not found",
			func() *types.MsgExecuteOrders {
				return &types.MsgExecuteOrders{
					Creator:           addr[2].String(),
					SpotOrderIds:      []uint64{1},
					PerpetualOrderIds: []uint64{1},
				}
			},
			func() {},
		},
		{
			"Perpetual order not found",
			"perpetual order not found",
			func() *types.MsgExecuteOrders {
				openOrderMsg := &types.MsgCreateSpotOrder{
					OwnerAddress: addr[2].String(),
					OrderType:    types.SpotOrderType_LIMITBUY,
					OrderPrice: types.OrderPrice{
						BaseDenom:  "uusdc",
						QuoteDenom: "uatom",
						Rate:       math.LegacyNewDec(5),
					},
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreateSpotOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				return &types.MsgExecuteOrders{
					Creator:           addr[2].String(),
					SpotOrderIds:      []uint64{1},
					PerpetualOrderIds: []uint64{1},
				}
			},
			func() {},
		},
		{
			"Success: Execute Spot Order",
			"",
			func() *types.MsgExecuteOrders {
				suite.SetupCoinPrices()

				_ = suite.CreateNewAmmPool(addr[0], true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, math.NewInt(100000000000).MulRaw(10), math.NewInt(100000000000).MulRaw(10))

				openOrderMsg := &types.MsgCreateSpotOrder{
					OwnerAddress: addr[2].String(),
					OrderType:    types.SpotOrderType_LIMITBUY,
					OrderPrice: types.OrderPrice{
						BaseDenom:  "uusdc",
						QuoteDenom: "uatom",
						Rate:       math.LegacyNewDec(10),
					},
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreateSpotOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(5),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgExecuteOrders{
					Creator:           addr[2].String(),
					SpotOrderIds:      []uint64{1},
					PerpetualOrderIds: []uint64{},
				}
			},
			func() {
				// Get events from context
				events := suite.ctx.EventManager().Events()

				// Find the specific event we're looking for
				var foundEvent sdk.Event
				for _, event := range events {
					if event.Type == types.TypeEvtExecuteLimitBuySpotOrder {
						foundEvent = event
						break
					}
				}

				// Assert event was emitted
				suite.Require().NotNil(foundEvent)

				// Check event attributes
				suite.Require().Equal(types.TypeEvtExecuteLimitBuySpotOrder, foundEvent.Type)

				// Check specific attributes
				for _, attr := range foundEvent.Attributes {
					switch string(attr.Key) {
					case "order_id":
						suite.Require().Equal("1", string(attr.Value))
					case "order_price":
						suite.Require().Equal(string(attr.Value), "{\"base_denom\":\"uusdc\",\"quote_denom\":\"uatom\",\"rate\":\"5.000000000000000000\"}")
					}
				}
			},
		},
		{
			"Success: Execute Perpetual Order",
			"",
			func() *types.MsgExecuteOrders {
				_, _, _ = suite.SetPerpetualPool(1)
				openOrderMsg := &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress: addr[2].String(),
					TriggerPrice: types.TriggerPrice{
						TradingAssetDenom: "uatom",
						Rate:              math.LegacyNewDec(10),
					},
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(100)},
					TradingAsset:    "uatom",
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(15),
					StopLossPrice:   math.LegacyNewDec(8),
					PoolId:          1,
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreatePerpetualOpenOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(10),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgExecuteOrders{
					Creator:           addr[2].String(),
					SpotOrderIds:      []uint64{},
					PerpetualOrderIds: []uint64{1},
				}
			},
			func() {
				// Get events from context
				events := suite.ctx.EventManager().Events()

				// Find the specific event we're looking for
				var foundEvent sdk.Event
				for _, event := range events {
					if event.Type == types.TypeEvtExecuteLimitOpenPerpetualOrder {
						foundEvent = event
						break
					}
				}

				// Assert event was emitted
				suite.Require().NotNil(foundEvent)

				// Check event attributes
				suite.Require().Equal(types.TypeEvtExecuteLimitOpenPerpetualOrder, foundEvent.Type)

				// Check specific attributes
				for _, attr := range foundEvent.Attributes {
					switch string(attr.Key) {
					case "order_id":
						suite.Require().Equal("1", string(attr.Value))
					case "trigger_price":
						suite.Require().Equal(string(attr.Value), "{\"trading_asset_denom\":\"uatom\",\"rate\":\"10.000000000000000000\"}")
					}
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.ExecuteOrders(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postrequisiteFunction()
		})
	}
}
