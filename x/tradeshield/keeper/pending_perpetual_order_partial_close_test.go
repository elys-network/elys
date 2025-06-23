package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/v6/x/oracle/types"
	"github.com/elys-network/elys/v6/x/tradeshield/keeper"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestExecutePartialLimitCloseOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                  string
		expectErrMsg          string
		prerequisiteFunction  func() *types.MsgExecuteOrders
		postrequisiteFunction func()
	}{
		{
			"Success: Execute Partial Limit Close Order (50%)",
			"",
			func() *types.MsgExecuteOrders {
				_, _, _ = suite.SetPerpetualPool(1)

				// First create and execute a limit open order to create a position
				openOrderMsg := &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(10),
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(1000)},
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(15),
					StopLossPrice:   math.LegacyNewDec(8),
					PoolId:          1,
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreatePerpetualOpenOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				// Execute the open order
				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(10),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				executeOpenMsg := &types.MsgExecuteOrders{
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
				_, err = msgSrvr.ExecuteOrders(suite.ctx, executeOpenMsg)
				suite.Require().NoError(err)

				// Now create a partial limit close order (50%)
				closeOrderMsg := &types.MsgCreatePerpetualCloseOrder{
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(12), // Higher price for LONG position
					PositionId:      1,
					ClosePercentage: 50, // 50% close
					PoolId:          1,
				}
				_, err = msgSrvr.CreatePerpetualCloseOrder(suite.ctx, closeOrderMsg)
				suite.Require().NoError(err)

				// Set price to trigger the close order
				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(12),
					Source:    "elys",
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
							OrderId:      2, // The close order
						},
					},
				}
			},
			func() {
				// Verify the close order was executed and removed
				_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(addr[2].String()), 1, 2)
				suite.Require().False(found, "Close order should be removed after execution")

				// Verify the position still exists but with reduced size
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, 1, sdk.MustAccAddressFromBech32(addr[2].String()), 1)
				suite.Require().NoError(err)
				suite.Require().NotNil(mtp)
				// The position should still exist but with reduced custody (50% of original)
				suite.Require().Equal(math.NewInt(2500), mtp.Custody)
			},
		},
		{
			"Success: Execute Partial Limit Close Order (25%)",
			"",
			func() *types.MsgExecuteOrders {
				suite.ResetSuite()
				addr = suite.AddAccounts(3, addr)
				_, _, _ = suite.SetPerpetualPool(1)

				// First create and execute a limit open order to create a position
				openOrderMsg := &types.MsgCreatePerpetualOpenOrder{
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(10),
					Collateral:      sdk.Coin{Denom: "uatom", Amount: math.NewInt(2000)},
					Position:        types.PerpetualPosition_LONG,
					Leverage:        math.LegacyNewDec(5),
					TakeProfitPrice: math.LegacyNewDec(15),
					StopLossPrice:   math.LegacyNewDec(8),
					PoolId:          1,
				}
				msgSrvr := keeper.NewMsgServerImpl(suite.app.TradeshieldKeeper)
				_, err := msgSrvr.CreatePerpetualOpenOrder(suite.ctx, openOrderMsg)
				suite.Require().NoError(err)

				// Execute the open order
				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(10),
					Source:    "elys",
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				executeOpenMsg := &types.MsgExecuteOrders{
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
				_, err = msgSrvr.ExecuteOrders(suite.ctx, executeOpenMsg)
				suite.Require().NoError(err)

				// Now create a partial limit close order (25%)
				closeOrderMsg := &types.MsgCreatePerpetualCloseOrder{
					OwnerAddress:    addr[2].String(),
					TriggerPrice:    math.LegacyNewDec(11), // Higher price for LONG position
					PositionId:      1,
					ClosePercentage: 25, // 25% close
					PoolId:          1,
				}
				_, err = msgSrvr.CreatePerpetualCloseOrder(suite.ctx, closeOrderMsg)
				suite.Require().NoError(err)

				// Set price to trigger the close order
				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(11),
					Source:    "elys",
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
							OrderId:      2, // The close order
						},
					},
				}
			},
			func() {
				// Verify the close order was executed and removed
				_, found := suite.app.TradeshieldKeeper.GetPendingPerpetualOrder(suite.ctx, sdk.MustAccAddressFromBech32(addr[2].String()), 1, 2)
				suite.Require().False(found, "Close order should be removed after execution")

				// Verify the position still exists but with reduced size
				mtp, err := suite.app.PerpetualKeeper.GetMTP(suite.ctx, 1, sdk.MustAccAddressFromBech32(addr[2].String()), 1)
				suite.Require().NoError(err)
				suite.Require().NotNil(mtp)
				// The position should still exist but with reduced custody (75% of original)
				suite.Require().Equal(math.NewInt(7500), mtp.Custody)
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
