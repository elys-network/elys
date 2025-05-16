package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (suite *PerpetualKeeperTestSuite) TestUpdateStopLossPrice() {

	// Define test cases
	testCases := []struct {
		name           string
		setup          func() *types.MsgUpdateStopLoss
		expectedErrMsg string
	}{
		{
			"mtp not found",
			func() *types.MsgUpdateStopLoss {
				return &types.MsgUpdateStopLoss{
					Creator: sample.AccAddress(),
					Id:      uint64(10),
					Price:   math.LegacyNewDec(2),
				}
			},
			"mtp not found",
		},
		{
			"perpetual pool does not exist",
			func() *types.MsgUpdateStopLoss {
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4).Dec(),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				suite.app.PerpetualKeeper.RemovePool(suite.ctx, ammPool.PoolId)
				return &types.MsgUpdateStopLoss{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(2),
				}
			},
			"perpetual pool does not exist",
		},
		{
			"asset profile not found",
			func() *types.MsgUpdateStopLoss {
				suite.ResetSuite()
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4).Dec(),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)
				suite.app.OracleKeeper.RemoveAssetInfo(suite.ctx, ptypes.ATOM)
				return &types.MsgUpdateStopLoss{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(2),
				}
			},
			"price for outToken not set: uatom",
		},
		{
			"success: Stop Loss price updated",
			func() *types.MsgUpdateStopLoss {
				suite.ResetSuite()
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4).Dec(),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				return &types.MsgUpdateStopLoss{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(4),
				}
			},
			"",
		},
		{
			"updated Stop Loss price is greater than current market price for long",
			func() *types.MsgUpdateStopLoss {
				suite.ResetSuite()
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, _, err := suite.app.PerpetualKeeper.GetAssetPriceAndAssetUsdcDenomRatio(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4).Dec(),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(25),
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgUpdateStopLoss{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(30),
				}
			},
			"stop loss price cannot be greater than equal to tradingAssetPrice for long",
		},
		{
			"updated Stop Loss price is less than current market price for short",
			func() *types.MsgUpdateStopLoss {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_SHORT,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("2.9"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyNewDec(3),
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgUpdateStopLoss{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyMustNewDecFromStr("2"),
				}
			},
			"stop loss price cannot be less than equal to tradingAssetPrice for short",
		},
		{
			"Force close when trying to update stop loss in unhealthy position",
			func() *types.MsgUpdateStopLoss {
				suite.ResetSuite()

				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyMustNewDecFromStr("2"),
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(5),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: math.LegacyMustNewDecFromStr("2.9"),
					StopLossPrice:   math.LegacyZeroDec(),
				}
				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg)
				suite.Require().NoError(err)

				suite.app.OracleKeeper.SetPrice(suite.ctx, oracletypes.Price{
					Asset:     "ATOM",
					Price:     math.LegacyMustNewDecFromStr("0.00001"),
					Provider:  oracleProvider.String(),
					Timestamp: uint64(suite.ctx.BlockTime().Unix()),
				})

				return &types.MsgUpdateStopLoss{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyMustNewDecFromStr("2"),
				}
			},
			"",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			msgSrvr := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
			_, err := msgSrvr.UpdateStopLoss(suite.ctx, msg)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
