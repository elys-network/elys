package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	simapp "github.com/elys-network/elys/v4/app"
	"github.com/elys-network/elys/v4/x/leveragelp/keeper"
	"github.com/elys-network/elys/v4/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

func (suite *KeeperTestSuite) TestCloseositions() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := sdkmath.LegacyNewDec(2)
	testCases := []struct {
		name                 string
		input                *types.MsgClosePositions
		prerequisiteFunction func() *types.Position
		postValidateFunction func()
	}{
		{
			"CheckAndLiquidateUnhealthyPosition returns error: token prices not set",
			&types.MsgClosePositions{
				Creator: addresses[7].String(),
				Liquidate: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      1,
					},
				},
				StopLoss: nil,
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyNewDec(2),
				}
				suite.SetCurrentHeight(1)
				params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				params.FallbackEnabled = true
				params.NumberPerBlock = 1
				err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				// doing this before gives panic
				suite.RemovePrices(suite.ctx, []string{"uusdc"})
				return position1
			},
			func() {

			},
		},
		{
			"CheckAndCloseAtStopLoss returns error: token prices not set",
			&types.MsgClosePositions{
				Creator:   addresses[7].String(),
				Liquidate: nil,
				StopLoss: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      1,
					},
				},
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyNewDec(2),
				}
				suite.SetCurrentHeight(1)
				params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				params.FallbackEnabled = true
				params.NumberPerBlock = 1
				err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				// doing this before gives panic
				suite.RemovePrices(suite.ctx, []string{"uusdc"})
				return position1
			},
			func() {

			},
		},
		{
			"pool not found in leveragelp for unhealthy position",
			&types.MsgClosePositions{
				Creator: addresses[7].String(),
				Liquidate: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      1,
					},
				},
				StopLoss: nil,
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				suite.app.LeveragelpKeeper.RemovePool(suite.ctx, 1)
				return position1
			},
			func() {
			},
		},
		{
			"pool not found in amm for unhealthy position",
			&types.MsgClosePositions{
				Creator: addresses[7].String(),
				Liquidate: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      1,
					},
				},
				StopLoss: nil,
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				suite.app.AmmKeeper.RemovePool(suite.ctx, 1)
				return position1
			},
			func() {
			},
		},
		{
			"pool not found in leveragelp for stop loss position",
			&types.MsgClosePositions{
				Creator:   addresses[7].String(),
				Liquidate: nil,
				StopLoss: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      1,
					},
				},
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				suite.app.LeveragelpKeeper.RemovePool(suite.ctx, 1)
				return position1
			},
			func() {
			},
		},
		{
			"pool not found in amm for stop loss position",
			&types.MsgClosePositions{
				Creator:   addresses[7].String(),
				Liquidate: nil,
				StopLoss: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      1,
					},
				},
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				suite.app.AmmKeeper.RemovePool(suite.ctx, 1)
				return position1
			},
			func() {
			},
		},
		{
			"position not found in leveragelp for unhealthy position",
			&types.MsgClosePositions{
				Creator: addresses[7].String(),
				Liquidate: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      2,
					},
				},
				StopLoss: nil,
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				return position1
			},
			func() {
			},
		},
		{
			"position not found in amm for stop loss position",
			&types.MsgClosePositions{
				Creator:   addresses[7].String(),
				Liquidate: nil,
				StopLoss: []*types.PositionRequest{
					{
						Address: addresses[1].String(),
						Id:      2,
					},
				},
			},
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: sdkmath.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1, 1)
				suite.Require().NoError(err)
				return position1
			},
			func() {
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_ = tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.ClosePositions(suite.ctx, tc.input)
			suite.Require().NoError(err)
			tc.postValidateFunction()
		})
	}
}
