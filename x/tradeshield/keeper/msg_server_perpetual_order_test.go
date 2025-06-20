package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/tradeshield/keeper"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
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
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(10),
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(100)},
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
					OwnerAddress:    addr[0].String(),
					TriggerPrice:    math.LegacyNewDec(10),
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(100)},
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
					PoolId:       1,
				}
			},
		},
		{
			"Incorrect Order Owner updating the order",
			"key 1 doesn't exist: key not found",
			func() *types.MsgUpdatePerpetualOrder {
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

				return &types.MsgUpdatePerpetualOrder{
					OwnerAddress: addr[1].String(), // incorrect owner
					OrderId:      1,
					PoolId:       1,
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
					PoolId:       1,
					TriggerPrice: math.LegacyNewDec(12),
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
			"order 1 doesn't exist: key not found",
			func() *types.MsgCancelPerpetualOrder {
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

				return &types.MsgCancelPerpetualOrder{
					OwnerAddress: addr[1].String(), // incorrect owner
					OrderId:      1,
					PoolId:       1,
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
					PoolId:       1,
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
				_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), 1, res.OrderId)
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
					Orders:       []*types.PerpetualOrderPoolKey{},
				}
			},
		},
		{
			"Success: close multiple orders",
			"",
			func() *types.MsgCancelPerpetualOrders {
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

				return &types.MsgCancelPerpetualOrders{
					OwnerAddress: addr[2].String(),
					Orders: []*types.PerpetualOrderPoolKey{
						{
							PoolId:  1,
							OrderId: 1,
						},
					},
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

func (suite *TradeshieldKeeperTestSuite) TestMsgServerCancelAllPerpetualOrders() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCancelAllPerpetualOrders
	}{
		{
			"No order for cancelling",
			"perpetual order not found",
			func() *types.MsgCancelAllPerpetualOrders {
				return &types.MsgCancelAllPerpetualOrders{
					OwnerAddress: addr[2].String(),
				}
			},
		},
		{
			"Success: close All orders",
			"",
			func() *types.MsgCancelAllPerpetualOrders {
				_, _, _ = suite.SetPerpetualPool(1)
				_, _, _ = suite.SetPerpetualPool(2)

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

				openOrderMsg2 := &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(5),
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(100)},
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(15),
					StopLossPrice:   math.LegacyNewDec(1),
					PoolId:          2,
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreatePerpetualOpenOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)
				_, err = msgSrvr.CreatePerpetualOpenOrder(suite.ctx, openOrderMsg2)
				suite.Require().NoError(err)

				return &types.MsgCancelAllPerpetualOrders{
					OwnerAddress: addr[2].String(),
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.CancelAllPerpetualOrders(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				status := types.Status_PENDING
				orders, _, err := suite.app.TradeshieldKeeper.GetPendingPerpetualOrdersForAddress(suite.ctx, msg.OwnerAddress, &status, nil)
				suite.Require().NoError(err)
				suite.Require().Empty(orders, "All orders should be cancelled")
			}
		})
	}
}
