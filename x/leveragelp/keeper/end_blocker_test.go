package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v7/app"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
	"time"
)

func (suite *KeeperTestSuite) TestEndBlocker() {
	suite.ResetSuite()
	suite.SetupCoinPrices(suite.ctx)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 5, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForADL(suite, addresses, asset1, asset2)

	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.NumberPerBlock = 2
	err := suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)

	leverage := sdkmath.LegacyMustNewDecFromStr("2")
	collateralAmount := sdkmath.NewInt(5_000_000_000)

	leveragePool, found := suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
	suite.Require().True(found)

	ammPool, err := suite.app.LeveragelpKeeper.GetAmmPool(suite.ctx, 1)
	suite.Require().NoError(err)
	currentLeverageRatio := leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	for _, address := range addresses {
		leveragePool, found = suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
		suite.Require().True(found)

		ammPool, err = suite.app.LeveragelpKeeper.GetAmmPool(suite.ctx, 1)
		suite.Require().NoError(err)
		currentLeverageRatio = leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

		if currentLeverageRatio.GT(leveragePool.MaxLeveragelpRatio) {
			break
		}
		msg := types.MsgOpen{
			Creator:          address.String(),
			CollateralAsset:  ptypes.BaseCurrency,
			CollateralAmount: collateralAmount,
			AmmPoolId:        1,
			Leverage:         leverage,
			StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("1.000000005"),
		}
		cacheCtx, write := suite.ctx.CacheContext()
		_, err = suite.app.LeveragelpKeeper.Open(cacheCtx, &msg)
		if err == nil {
			write()
		}
	}

	suite.ctx = suite.ctx.WithBlockHeight(suite.ctx.BlockHeight() + 10000).WithBlockTime(suite.ctx.BlockTime().Add(2 * time.Hour))

	exitPool := ammtypes.MsgExitPool{
		Sender:        addresses[0].String(),
		PoolId:        1,
		MinAmountsOut: sdk.NewCoins(sdk.NewInt64Coin(ptypes.ATOM, 1), sdk.NewInt64Coin(ptypes.BaseCurrency, 1)),
		ShareAmountIn: sdkmath.LegacyMustNewDecFromStr("300000000000000000000000").TruncateInt(),
		TokenOutDenom: "",
	}

	_, _, _, _, _, err = suite.app.AmmKeeper.ExitPool(suite.ctx, addresses[0], 1, exitPool.ShareAmountIn, exitPool.MinAmountsOut, exitPool.TokenOutDenom, false, false)
	suite.Require().NoError(err)

	leveragePool, found = suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
	suite.Require().True(found)

	ammPool, err = suite.app.LeveragelpKeeper.GetAmmPool(suite.ctx, 1)
	suite.Require().NoError(err)
	currentLeverageRatio = leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	suite.app.LeveragelpKeeper.EndBlocker(suite.ctx)

	leveragePool, found = suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
	suite.Require().True(found)

	ammPool, err = suite.app.LeveragelpKeeper.GetAmmPool(suite.ctx, 1)
	suite.Require().NoError(err)
	currentLeverageRatio = leveragePool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

	suite.app.LeveragelpKeeper.EndBlocker(suite.ctx)
}
