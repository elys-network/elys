package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	keeper "github.com/elys-network/elys/x/tradeshield/keeper"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (suite *TradeshieldKeeperTestSuite) TestMsgServerExecuteOrder() {
	addr := suite.AddAccounts(3, nil)

	testCases := []struct {
		name                 string
		expectErrMsg         string
		prerequisiteFunction func() *types.MsgExecuteOrders
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
		},
		{
			"Success: Execute Orders",
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
					SpotOrderIds:      []uint64{1},
					PerpetualOrderIds: []uint64{1},
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
		})
	}
}
