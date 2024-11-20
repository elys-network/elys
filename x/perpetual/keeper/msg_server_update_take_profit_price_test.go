package keeper_test

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/testutil/sample"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/elys-network/elys/x/perpetual/keeper"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (suite *PerpetualKeeperTestSuite) TestUpdateTakeProfitPrice() {

	// Define test cases
	testCases := []struct {
		name           string
		setup          func() *types.MsgUpdateTakeProfitPrice
		expectedErrMsg string
	}{
		{
			"mtp not found",
			func() *types.MsgUpdateTakeProfitPrice {
				return &types.MsgUpdateTakeProfitPrice{
					Creator: sample.AccAddress(),
					Id:      uint64(10),
					Price:   math.LegacyNewDec(2),
				}
			},
			"mtp not found",
		},
		{
			"perpetual pool does not exist",
			func() *types.MsgUpdateTakeProfitPrice {
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
				suite.Require().NoError(err)
				suite.app.PerpetualKeeper.RemovePool(suite.ctx, ammPool.PoolId)
				return &types.MsgUpdateTakeProfitPrice{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(2),
				}
			},
			"perpetual pool does not exist",
		},
		{
			"asset profile not found",
			func() *types.MsgUpdateTakeProfitPrice {
				suite.ResetSuite()
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
				suite.Require().NoError(err)
				suite.app.OracleKeeper.RemoveAssetInfo(suite.ctx, ptypes.ATOM)
				return &types.MsgUpdateTakeProfitPrice{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(2),
				}
			},
			"asset price uatom not found",
		},
		{
			"success: take profit price updated",
			func() *types.MsgUpdateTakeProfitPrice {
				suite.ResetSuite()
				addr := suite.AddAccounts(1, nil)
				positionCreator := addr[0]
				_, _, ammPool := suite.SetPerpetualPool(1)
				tradingAssetPrice, err := suite.app.PerpetualKeeper.GetAssetPrice(suite.ctx, ptypes.ATOM)
				suite.Require().NoError(err)
				openPositionMsg := &types.MsgOpen{
					Creator:         positionCreator.String(),
					Leverage:        math.LegacyNewDec(2),
					Position:        types.Position_LONG,
					PoolId:          ammPool.PoolId,
					TradingAsset:    ptypes.ATOM,
					Collateral:      sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(1000)),
					TakeProfitPrice: tradingAssetPrice.MulInt64(4),
					StopLossPrice:   math.LegacyZeroDec(),
				}

				position, err := suite.app.PerpetualKeeper.Open(suite.ctx, openPositionMsg, false)
				suite.Require().NoError(err)

				return &types.MsgUpdateTakeProfitPrice{
					Creator: positionCreator.String(),
					Id:      position.Id,
					Price:   math.LegacyNewDec(10),
				}
			},
			"",
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			msg := tc.setup()
			msgSrvr := keeper.NewMsgServerImpl(*suite.app.PerpetualKeeper)
			_, err := msgSrvr.UpdateTakeProfitPrice(suite.ctx, msg)

			if tc.expectedErrMsg != "" {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectedErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
