package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	keeper "github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestMsgServerCreateSpotOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCreateSpotOrder
		checkBalance         bool
	}{
		{
			"Success: Spot-Market Buy Order",
			"",
			func() *types.MsgCreateSpotOrder {
				// sets amm pool with prices
				_, _, _ = suite.SetPerpetualPool(1)
				return &types.MsgCreateSpotOrder{
					OwnerAddress:     addr[2].String(),
					OrderType:        types.SpotOrderType_MARKETBUY,
					OrderPrice:       math.LegacyNewDec(5),
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
			},
			false,
		},
		{
			"Success: Create Spot Order",
			"",
			func() *types.MsgCreateSpotOrder {
				return &types.MsgCreateSpotOrder{
					OwnerAddress:     addr[2].String(),
					OrderType:        types.SpotOrderType_LIMITBUY,
					OrderPrice:       math.LegacyNewDec(5),
					OrderAmount:      sdk.NewCoin("uusdc", math.NewInt(100000)),
					OrderTargetDenom: "uatom",
				}
			},
			true,
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			balanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uusdc")
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.CreateSpotOrder(suite.ctx, msg)
			balanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uusdc")
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				if tc.checkBalance {
					suite.Require().Equal(balanceBefore.Amount.Sub(msg.OrderAmount.Amount), balanceAfter.Amount)
				}
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServerUpdateSpotOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgUpdateSpotOrder
	}{
		{
			"Order not found",
			"spot order not found",
			func() *types.MsgUpdateSpotOrder {
				return &types.MsgUpdateSpotOrder{
					OwnerAddress: addr[2].String(),
					OrderId:      1,
				}
			},
		},
		{
			"Incorrect Order-Owner updating the order",
			"incorrect owner",
			func() *types.MsgUpdateSpotOrder {
				_, _, _ = suite.SetPerpetualPool(1)

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

				return &types.MsgUpdateSpotOrder{
					OwnerAddress: addr[1].String(), // incorrect owner
					OrderId:      1,
				}
			},
		},
		{
			"Success: Update Spot Order",
			"",
			func() *types.MsgUpdateSpotOrder {
				return &types.MsgUpdateSpotOrder{
					OwnerAddress: addr[2].String(),
					OrderId:      1,
					OrderPrice:   math.LegacyNewDec(12),
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.UpdateSpotOrder(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServerCancelSpotOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCancelSpotOrder
	}{
		{
			"Order not found",
			"spot order not found",
			func() *types.MsgCancelSpotOrder {
				return &types.MsgCancelSpotOrder{
					OwnerAddress: addr[2].String(),
					OrderId:      1,
				}
			},
		},
		{
			"Incorrect Order Owner cancelling the order",
			"incorrect owner",
			func() *types.MsgCancelSpotOrder {
				_, _, _ = suite.SetPerpetualPool(1)

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

				return &types.MsgCancelSpotOrder{
					OwnerAddress: addr[1].String(), // incorrect owner
					OrderId:      1,
				}
			},
		},
		{
			"Success: Cancel Spot Open Order",
			"",
			func() *types.MsgCancelSpotOrder {

				return &types.MsgCancelSpotOrder{
					OwnerAddress: addr[2].String(), // incorrect owner
					OrderId:      1,
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			balanceBefore := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uusdc")
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			res, err := msgSrvr.CancelSpotOrder(suite.ctx, msg)
			balanceAfter := suite.app.BankKeeper.GetBalance(suite.ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), "uusdc")
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
				_, found := suite.app.TradeshieldKeeper.GetPendingSpotOrder(suite.ctx, res.OrderId)
				suite.Require().False(found)
				// Hardcoded collateral amount, as using only one test case
				suite.Require().Equal(balanceBefore.Amount.Add(math.NewInt(100000)), balanceAfter.Amount)
			}
		})
	}
}

func (suite *TradeshieldKeeperTestSuite) TestMsgServerCancelSpotOrders() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgCancelSpotOrders
	}{
		{
			"No Ids sent for cancelling",
			"zero order ids",
			func() *types.MsgCancelSpotOrders {
				return &types.MsgCancelSpotOrders{
					Creator:      addr[2].String(),
					SpotOrderIds: []uint64{},
				}
			},
		},
		{
			"Success: close multiple orders",
			"",
			func() *types.MsgCancelSpotOrders {
				_, _, _ = suite.SetPerpetualPool(1)

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

				return &types.MsgCancelSpotOrders{
					Creator:      addr[2].String(),
					SpotOrderIds: []uint64{1},
				}
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.prerequisiteFunction()
			msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
			_, err := msgSrvr.CancelSpotOrders(suite.ctx, msg)
			if tc.expectErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
