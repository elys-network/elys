package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"github.com/elys-network/elys/v7/x/tradeshield/keeper"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (suite *TradeshieldKeeperTestSuite) TestMsgServerExecuteOrder() {
	addr := suite.AddAccounts(3, nil)

	perpParams := suite.app.PerpetualKeeper.GetParams(suite.ctx)
	perpParams.EnabledPools = []uint64{1}
	err := suite.app.PerpetualKeeper.SetParams(suite.ctx, &perpParams)
	suite.Require().NoError(err)

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
					Creator:      addr[2].String(),
					SpotOrderIds: []uint64{1},
					PerpetualOrders: []types.PerpetualOrderKey{
						{
							OwnerAddress: addr[2].String(),
							PoolId:       1,
							OrderId:      1,
						},
					},
				}
			},
			func() {},
		},
		{
			"Perpetual order not found",
			"perpetual order not found",
			func() *types.MsgExecuteOrders {
				openOrderMsg := &types.MsgCreateSpotOrder{
					OwnerAddress:     addr[2].String(),
					OrderType:        types.SpotOrderType_LIMITBUY,
					OrderPrice:       math.LegacyNewDec(5),
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreateSpotOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				return &types.MsgExecuteOrders{
					Creator:      addr[2].String(),
					SpotOrderIds: []uint64{1},
					PerpetualOrders: []types.PerpetualOrderKey{
						{
							OwnerAddress: addr[2].String(),
							PoolId:       1,
							OrderId:      1,
						},
					},
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
					OwnerAddress:     addr[2].String(),
					OrderType:        types.SpotOrderType_LIMITBUY,
					OrderPrice:       math.LegacyNewDec(10),
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreateSpotOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(5),
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgExecuteOrders{
					Creator:         addr[2].String(),
					SpotOrderIds:    []uint64{1},
					PerpetualOrders: []types.PerpetualOrderKey{},
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
						suite.Require().Equal(string(attr.Value), "\"5.000000000000000000\"")
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
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(10),
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(100)},
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
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgExecuteOrders{
					Creator:      addr[2].String(),
					SpotOrderIds: []uint64{},
					PerpetualOrders: []types.PerpetualOrderKey{
						{
							OwnerAddress: addr[2].String(),
							PoolId:       1,
							OrderId:      1,
						},
					},
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
						suite.Require().Equal(string(attr.Value), "\"10.000000000000000000\"")
					}
				}
			},
		},
		{
			"Success: Execute Multiple Spot Orders with One Failure",
			"",
			func() *types.MsgExecuteOrders {
				suite.ResetSuite()
				suite.SetupCoinPrices()

				addr = suite.AddAccounts(1, addr)
				_ = suite.CreateNewAmmPool(addr[0], true, math.LegacyZeroDec(), math.LegacyZeroDec(), ptypes.ATOM, math.NewInt(100000000000).MulRaw(10), math.NewInt(100000000000).MulRaw(10))

				openOrderMsg1 := &types.MsgCreateSpotOrder{
					OwnerAddress:     addr[2].String(),
					OrderType:        types.SpotOrderType_LIMITBUY,
					OrderPrice:       math.LegacyNewDec(10),
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreateSpotOrder(suite.ctx, openOrderMsg1)
				suite.Require().NoError(err)

				openOrderMsg2 := &types.MsgCreateSpotOrder{
					OwnerAddress:     addr[2].String(),
					OrderType:        types.SpotOrderType_LIMITSELL,
					OrderPrice:       math.LegacyNewDec(10),
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(200000)),
					OrderTargetDenom: "uatom",
				}
				_, err = msgSrvr.CreateSpotOrder(suite.ctx, openOrderMsg2)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(5),
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				// Return message with both order IDs
				return &types.MsgExecuteOrders{
					Creator:         addr[2].String(),
					SpotOrderIds:    []uint64{1, 2}, // Both orders exist but second will fail during execution
					PerpetualOrders: []types.PerpetualOrderKey{},
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

				// Assert event was emitted for the successful order
				suite.Require().NotNil(foundEvent)

				// Check event attributes
				suite.Require().Equal(types.TypeEvtExecuteLimitBuySpotOrder, foundEvent.Type)

				// Check specific attributes
				for _, attr := range foundEvent.Attributes {
					switch string(attr.Key) {
					case "order_id":
						suite.Require().Equal("1", string(attr.Value))
					case "order_price":
						suite.Require().Equal(string(attr.Value), "\"10.000000000000000000\"")
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
