package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	keeper "github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// TODO: Add test for CreatePerpetualCloseOrder after enabling the code
func (suite *TradeshieldKeeperTestSuite) TestMsgServerPerpetualOpenOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCreatePerpetualOpenOrder
	}{
		{
			"Pool not found",
			"pool 1 not found",
			func() *types.MsgCreatePerpetualOpenOrder {
				return &types.MsgCreatePerpetualOpenOrder{
					PoolId:       1,
					OwnerAddress: addr[2].String(),
				}
			},
		},
		{
			"Success: Create Perpetual Open Order",
			"",
			func() *types.MsgCreatePerpetualOpenOrder {
				_, _, _ = suite.SetPerpetualPool(1)
				return &types.MsgCreatePerpetualOpenOrder{
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
			},
		},
		{
			"Perpetual Open Order already pending for the same Pool", // From above test case
			"user already has a order for the same pool",
			func() *types.MsgCreatePerpetualOpenOrder {

				return &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress: addr[2].String(),
					TriggerPrice: types.TriggerPrice{
						TradingAssetDenom: "uatom",
						Rate:              math.LegacyNewDec(10),
					},
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(200)},
					TradingAsset:    "uatom",
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(15),
					StopLossPrice:   math.LegacyNewDec(8),
					PoolId:          1,
				}
			},
		},
		{
			"Position already open in the perpetual pool",
			"user already has a position in the same pool",
			func() *types.MsgCreatePerpetualOpenOrder {
				// Make asset price equal to the trigger price
				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(10),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.ExecuteOrders(suite.ctx, &types.MsgExecuteOrders{
					Creator:           addr[2].String(),
					PerpetualOrderIds: []uint64{1},
					SpotOrderIds:      []uint64{},
				})
				suite.T().Log("Error: ", err)
				suite.Require().NoError(err)

				return &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress: addr[2].String(),
					TriggerPrice: types.TriggerPrice{
						TradingAssetDenom: "uatom",
						Rate:              math.LegacyNewDec(10),
					},
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(200)},
					TradingAsset:    "uatom",
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(15),
					StopLossPrice:   math.LegacyNewDec(8),
					PoolId:          1,
				}
			},
		},
		{
			"Invalid PerpetualOrder, takeprofitprice is less than trigger price",
			"take profit price cannot be less than equal to trading price for long",
			func() *types.MsgCreatePerpetualOpenOrder {
				return &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress: addr[0].String(),
					TriggerPrice: types.TriggerPrice{
						TradingAssetDenom: "uatom",
						Rate:              math.LegacyNewDec(10),
					},
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(100)},
					TradingAsset:    "uatom",
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(2),
					StopLossPrice:   math.LegacyNewDec(8),
					PoolId:          1,
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			balanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uatom")
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.CreatePerpetualOpenOrder(suite.ctx, msg)
			balanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uatom")
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				suite.Require().Equal(balanceBefore.Amount.Sub(msg.Collateral.Amount), balanceAfter.Amount)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServerUpdatePerpetualOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgUpdatePerpetualOrder
	}{
		{
			"Order not found",
			"key 1 doesn't exist",
			func() *types.MsgUpdatePerpetualOrder {
				return &types.MsgUpdatePerpetualOrder{
					OwnerAddress: addr[2].String(),
					OrderId:      1,
				}
			},
		},
		{
			"Incorrect Order Owner updating the order",
			"incorrect owner",
			func() *types.MsgUpdatePerpetualOrder {
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

				return &types.MsgUpdatePerpetualOrder{
					OwnerAddress: addr[1].String(), // incorrect owner
					OrderId:      1,
				}
			},
		},
		// TODO: Add checks in code for triggerprice comparison with TakeProfitPrice and StopLossPrice
		// Update Rate with >15 should not pass, but passing here
		{
			"Success: Update Perpetual Open Order",
			"",
			func() *types.MsgUpdatePerpetualOrder {
				return &types.MsgUpdatePerpetualOrder{
					OwnerAddress: addr[2].String(),
					OrderId:      1,
					TriggerPrice: types.TriggerPrice{
						TradingAssetDenom: "uatom",
						Rate:              math.LegacyNewDec(12),
					},
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.UpdatePerpetualOrder(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServerCancelPerpetualOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCancelPerpetualOrder
	}{
		{
			"Order not found",
			"order 1 doesn't exist",
			func() *types.MsgCancelPerpetualOrder {
				return &types.MsgCancelPerpetualOrder{
					OwnerAddress: addr[2].String(),
					OrderId:      1,
				}
			},
		},
		{
			"Incorrect Order Owner cancelling the order",
			"incorrect owner",
			func() *types.MsgCancelPerpetualOrder {
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

				return &types.MsgCancelPerpetualOrder{
					OwnerAddress: addr[1].String(), // incorrect owner
					OrderId:      1,
				}
			},
		},
		{
			"Success: Cancel Perpetual Open Order",
			"",
			func() *types.MsgCancelPerpetualOrder {

				return &types.MsgCancelPerpetualOrder{
					OwnerAddress: addr[2].String(), // incorrect owner
					OrderId:      1,
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			balanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uatom")
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			res, err := msgSrvr.CancelPerpetualOrder(suite.ctx, msg)
			balanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uatom")
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, res.OrderId)
				suite.Require().False(found)
				// Hardcoded collateral amount, as using only one test case
				suite.Require().Equal(balanceBefore.Amount.Add(math.NewInt(100)), balanceAfter.Amount)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServerCancelPerpetualOrders() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCancelPerpetualOrders
	}{
		{
			"No Ids sent for cancelling",
			"zero order ids",
			func() *types.MsgCancelPerpetualOrders {
				return &types.MsgCancelPerpetualOrders{
					OwnerAddress: addr[2].String(),
					OrderIds:     []uint64{},
				}
			},
		},
		{
			"Success: close multiple orders",
			"",
			func() *types.MsgCancelPerpetualOrders {
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

				return &types.MsgCancelPerpetualOrders{
					OwnerAddress: addr[2].String(),
					OrderIds:     []uint64{1},
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.CancelPerpetualOrders(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
