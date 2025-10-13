package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/v6/app"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/leveragelp/keeper"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/v6/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/v6/x/stablestake/types"
	"time"
)

func (suite *KeeperTestSuite) TestGetSetADLCounter() {
	val := types.ADLCounter{
		PoolId:  1,
		Counter: 5,
	}
	suite.app.LeveragelpKeeper.SetADLCounter(suite.ctx, val)
	got := suite.app.LeveragelpKeeper.GetADLCounter(suite.ctx, val.PoolId)
	suite.Require().Equal(val, got)

	val.PoolId = 7
	val.Counter = 9
	suite.app.LeveragelpKeeper.SetADLCounter(suite.ctx, val)
	got = suite.app.LeveragelpKeeper.GetADLCounter(suite.ctx, val.PoolId)
	suite.Require().Equal(val, got)

	all := suite.app.LeveragelpKeeper.GetAllADLCounter(suite.ctx)
	suite.Require().Len(all, 2)

}

func initializeForADL(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string) {
	fee := sdkmath.LegacyMustNewDecFromStr("0.0002")
	issueAmount := sdkmath.NewInt(200_000_000_000)
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount.MulRaw(1000)),
			sdk.NewCoin(ptypes.Elys, issueAmount.MulRaw(1000)),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount.MulRaw(1000)),
		)
		err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
		if err != nil {
			panic(err)
		}
		err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, address, coins)
		if err != nil {
			panic(err)
		}
	}
	msgCreatePool := ammtypes.MsgCreatePool{
		Sender: addresses[0].String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:   fee,
			UseOracle: true,
			FeeDenom:  ptypes.Elys,
		},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewCoin(asset1, issueAmount.QuoRaw(2)),
				Weight: sdkmath.NewInt(50),
			},
			{
				Token:  sdk.NewCoin(asset2, issueAmount.QuoRaw(2)),
				Weight: sdkmath.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	suite.Require().NoError(err)
	enablePoolMsg := types.MsgAddPool{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Pool: types.AddPool{
			AmmPoolId:            poolId,
			LeverageMax:          sdkmath.LegacyNewDec(10),
			PoolMaxLeverageRatio: sdkmath.LegacyMustNewDecFromStr("0.05"),
			AdlTriggerRatio:      sdkmath.LegacyMustNewDecFromStr("0.06"),
		},
	}
	msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
	_, err = msgServer.AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(poolId, sdkmath.LegacyMustNewDecFromStr("10"), sdkmath.LegacyMustNewDecFromStr("0.05"), sdkmath.LegacyMustNewDecFromStr("0.06")))
	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  issueAmount.QuoRaw(2),
		PoolId:  1,
	}

	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.EnabledPools = []uint64{1}
	err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)

	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) TestClosePositionsOnADL() {
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

	err = suite.app.LeveragelpKeeper.ClosePositionsOnADL(suite.ctx, leveragePool)
	suite.Require().NoError(err)

	leveragePool, found = suite.app.LeveragelpKeeper.GetPool(suite.ctx, 1)
	suite.Require().True(found)

	err = suite.app.LeveragelpKeeper.ClosePositionsOnADL(suite.ctx, leveragePool)
	suite.Require().NoError(err)

}
