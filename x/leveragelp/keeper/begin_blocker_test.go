package keeper_test

import (
	"cosmossdk.io/math"
	"time"

	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (suite *KeeperTestSuite) TestBeginBlocker() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, math.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := math.LegacyNewDec(2)
	testCases := []struct {
		name                 string
		prerequisiteFunction func() *types.Position
		postValidateFunction func()
	}{
		{
			"CheckAndLiquidateUnhealthyPosition returns error: token prices not set",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				openMsg2 := types.MsgOpen{
					Creator:          addresses[2].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				openMsg3 := types.MsgOpen{
					Creator:          addresses[3].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				suite.SetCurrentHeight(1)
				params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				params.FallbackEnabled = true
				params.NumberPerBlock = 1
				err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1)
				suite.Require().NoError(err)
				_, err = suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg2)
				suite.Require().NoError(err)
				_, err = suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg3)
				suite.Require().NoError(err)
				// doing this before gives panic
				suite.RemovePrices(suite.ctx, []string{"uusdc"})
				return position1
			},
			func() {

			},
		},
		{
			"multiple closing positions in same amm pool, one is for stop loss, others are for low health",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyMustNewDecFromStr("1.5"),
				}
				openMsg2 := types.MsgOpen{
					Creator:          addresses[2].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage.MulInt64(2),
					StopLossPrice:    math.LegacyNewDec(2),
				}
				openMsg3 := types.MsgOpen{
					Creator:          addresses[3].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage.MulInt64(5),
					StopLossPrice:    math.LegacyNewDec(2),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1)
				suite.Require().NoError(err)
				_, err = suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg2)
				suite.Require().NoError(err)
				_, err = suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg3)
				suite.Require().NoError(err)
				params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				params.FallbackEnabled = true
				params.NumberPerBlock = 5
				params.SafetyFactor = math.LegacyMustNewDecFromStr("1.75")
				err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				suite.AddBlockTime(2 * time.Hour)
				return position1
			},
			func() {
				offset, _ := suite.app.LeveragelpKeeper.GetOffset(suite.ctx)
				suite.Require().Equal(offset, uint64(0))
			},
		},
		{
			"pool not found in leveragelp",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1)
				suite.Require().NoError(err)
				suite.app.LeveragelpKeeper.RemovePool(suite.ctx, 1)
				return position1
			},
			func() {
			},
		},
		{
			"pool not found in amm",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg1 := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyMustNewDecFromStr("1.5"),
				}
				suite.SetCurrentHeight(1)
				initializeForOpen(suite, addresses, asset1, asset2)
				position1, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg1)
				suite.Require().NoError(err)
				suite.app.AmmKeeper.RemovePool(suite.ctx, 1)
				return position1
			},
			func() {
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			_ = tc.prerequisiteFunction()
			suite.app.LeveragelpKeeper.BeginBlocker(suite.ctx)
			tc.postValidateFunction()
		})
	}
}

func (suite *KeeperTestSuite) TestCheckAndLiquidateUnhealthyPosition() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, math.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := math.LegacyNewDec(5)
	testCases := []struct {
		name                 string
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func() *types.Position
		postValidateFunction func(isHealthy, closeAttempted bool)
	}{
		{
			"token prices not set",
			true,
			"token price not set: uusdc",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				initializeForOpen(suite, addresses, asset1, asset2)
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				// doing this before gives panic
				suite.RemovePrices(suite.ctx, []string{"uusdc"})
				return position
			},
			func(isHealthy, closeAttempted bool) {
				suite.Require().False(isHealthy)
				suite.Require().False(closeAttempted)
			},
		},
		{
			"position is healthy",
			true,
			"position is healthy to close",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				return position
			},
			func(isHealthy, closeAttempted bool) {
				suite.Require().True(isHealthy)
				suite.Require().False(closeAttempted)
			},
		},
		{
			"success", // liquidating position before lockup period should be successful
			false,
			"",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				params.SafetyFactor = math.LegacyOneDec().MulInt64(1000)
				err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				return position
			},
			func(isHealthy, closeAttempted bool) {
				suite.Require().False(isHealthy)
				suite.Require().True(closeAttempted)
			},
		},
		{
			"success",
			false,
			"",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
				params.SafetyFactor = math.LegacyOneDec().MulInt64(1000)
				err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
				suite.Require().NoError(err)
				suite.AddBlockTime(2 * time.Hour)
				return position
			},
			func(isHealthy, closeAttempted bool) {
				suite.Require().False(isHealthy)
				suite.Require().True(closeAttempted)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			position := tc.prerequisiteFunction()
			pool, found := suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
			suite.Require().True(found)

			isHealthy, closeAttempted, _, err := suite.app.LeveragelpKeeper.CheckAndLiquidateUnhealthyPosition(suite.ctx, position, pool)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunction(isHealthy, closeAttempted)
		})
	}

}

func (suite *KeeperTestSuite) TestCheckAndCloseAtStopLoss() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, math.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := math.LegacyNewDec(5)
	testCases := []struct {
		name                 string
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func() *types.Position
		postValidateFunction func(underStopLossPrice, closeAttempted bool)
	}{
		{
			"token prices not set",
			true,
			"token price not set: uusdc",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				initializeForOpen(suite, addresses, asset1, asset2)
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				// doing this before gives panic
				suite.RemovePrices(suite.ctx, []string{"uusdc"})
				return position
			},
			func(underStopLossPrice, closeAttempted bool) {
				suite.Require().False(underStopLossPrice)
				suite.Require().False(closeAttempted)
			},
		},
		{
			"stop loss price > current price",
			true,
			"position stop loss price is not <= lp token price",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(1).QuoInt64(2),
				}
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				suite.AddBlockTime(2 * time.Hour)
				return position
			},
			func(underStopLossPrice, closeAttempted bool) {
				suite.Require().False(underStopLossPrice)
				suite.Require().False(closeAttempted)
			},
		},
		{
			"closing position before 1 hour",
			true,
			"funds will be locked for 1 hour",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				return position
			},
			func(underStopLossPrice, closeAttempted bool) {
				suite.Require().True(underStopLossPrice)
				suite.Require().True(closeAttempted)
			},
		},
		{
			"success",
			false,
			"",
			func() *types.Position {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForOpen(suite, addresses, asset1, asset2)
				openMsg := types.MsgOpen{
					Creator:          addresses[1].String(),
					CollateralAsset:  "uusdc",
					CollateralAmount: math.NewInt(1000_000),
					AmmPoolId:        1,
					Leverage:         leverage,
					StopLossPrice:    math.LegacyNewDec(2),
				}
				position, err := suite.app.LeveragelpKeeper.OpenLong(suite.ctx, &openMsg)
				suite.Require().NoError(err)
				suite.AddBlockTime(2 * time.Hour)
				return position
			},
			func(underStopLossPrice, closeAttempted bool) {
				suite.Require().True(underStopLossPrice)
				suite.Require().True(closeAttempted)
			},
		},
	}
	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			position := tc.prerequisiteFunction()
			pool, found := suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
			suite.Require().True(found)
			ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, 1)
			suite.Require().True(found)

			underStopLossPrice, closeAttempted, err := suite.app.LeveragelpKeeper.CheckAndCloseAtStopLoss(suite.ctx, position, pool, ammPool)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
			tc.postValidateFunction(underStopLossPrice, closeAttempted)
		})
	}

}
