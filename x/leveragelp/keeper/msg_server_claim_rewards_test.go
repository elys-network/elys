package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/v6/app"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/leveragelp/keeper"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	mastercheftypes "github.com/elys-network/elys/v6/x/masterchef/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/v6/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/v6/x/stablestake/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func initializeForClaimRewards(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string, createAmmPool bool) {
	fee := sdkmath.LegacyMustNewDecFromStr("0.0002")
	issueAmount := sdkmath.NewInt(10_000_000_000_000)
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount.MulRaw(100)),
			sdk.NewCoin(ptypes.Elys, issueAmount.MulRaw(100)),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount.MulRaw(100)),
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
				Token:  sdk.NewCoin(asset1, issueAmount),
				Weight: sdkmath.NewInt(50),
			},
			{
				Token:  sdk.NewCoin(asset2, issueAmount),
				Weight: sdkmath.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	suite.Require().NoError(err)
	if createAmmPool {
		enablePoolMsg := types.MsgAddPool{
			Authority: authtypes.NewModuleAddress("gov").String(),
			Pool: types.AddPool{
				AmmPoolId:            poolId,
				LeverageMax:          sdkmath.LegacyNewDec(10),
				PoolMaxLeverageRatio: sdkmath.LegacyMustNewDecFromStr("0.99"),
				AdlTriggerRatio:      sdkmath.LegacyMustNewDecFromStr("0.9999"),
			},
		}
		msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
		_, err = msgServer.AddPool(suite.ctx, &enablePoolMsg)
		suite.Require().NoError(err)

	}

	params := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	params.EnabledPools = []uint64{1}
	err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &params)
	suite.Require().NoError(err)

	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  issueAmount.QuoRaw(20),
		PoolId:  1,
	}
	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
	msgBond.Creator = addresses[2].String()
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
}

func openPosition(suite *KeeperTestSuite, address sdk.AccAddress, collateralAmount sdkmath.Int, leverage osmomath.BigDec) {
	msg := types.MsgOpen{
		Creator:          address.String(),
		CollateralAsset:  ptypes.BaseCurrency,
		CollateralAmount: collateralAmount,
		AmmPoolId:        1,
		Leverage:         leverage.Dec(),
		StopLossPrice:    sdkmath.LegacyMustNewDecFromStr("50.0"),
	}
	_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) TestMsgServerClaimRewards() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := osmomath.MustNewBigDecFromStr("2.0")
	collateralAmount := sdkmath.NewInt(10000000)
	testCases := []struct {
		name                 string
		input                *types.MsgClaimRewards
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"position not found",
			&types.MsgClaimRewards{
				Sender:     addresses[0].String(),
				PositionId: 1,
				PoolId:     1,
			},
			true,
			types.ErrPositionDoesNotExist.Error(),
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClaimRewards(suite, addresses, asset1, asset2, false)
			},
		},
		{"module is out of funds",
			&types.MsgClaimRewards{
				Sender:     addresses[0].String(),
				PositionId: 1,
				PoolId:     1,
			},
			true,
			"insufficient funds",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClaimRewards(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
				moduleAddress := address.Module("masterchef")
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, moduleAddress)
				err := suite.app.BankKeeper.SendCoins(suite.ctx, moduleAddress, addresses[2], balances)
				if err != nil {
					panic(err)
				}
				positonAddress := types.GetPositionAddress(1)
				suite.app.MasterchefKeeper.SetUserRewardInfo(suite.ctx, mastercheftypes.UserRewardInfo{
					User:          positonAddress.String(),
					PoolId:        1,
					RewardDenom:   "uusdc",
					RewardPending: sdkmath.LegacyMustNewDecFromStr("100"),
				})
			},
		},
		{"positive case",
			&types.MsgClaimRewards{
				Sender:     addresses[0].String(),
				PositionId: 1,
				PoolId:     1,
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClaimRewards(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.ClaimRewards(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}

func (suite *KeeperTestSuite) TestMsgServerClaimAllRewards() {
	suite.ResetSuite()
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdkmath.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := osmomath.MustNewBigDecFromStr("2.0")
	collateralAmount := sdkmath.NewInt(10000000)
	testCases := []struct {
		name                 string
		input                *types.MsgClaimAllUserRewards
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"module is out of funds",
			&types.MsgClaimAllUserRewards{
				Sender: addresses[0].String(),
			},
			true,
			"insufficient funds",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClaimRewards(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
				moduleAddress := address.Module("masterchef")
				balances := suite.app.BankKeeper.GetAllBalances(suite.ctx, moduleAddress)
				err := suite.app.BankKeeper.SendCoins(suite.ctx, moduleAddress, addresses[2], balances)
				if err != nil {
					panic(err)
				}
				positonAddress := types.GetPositionAddress(1)
				suite.app.MasterchefKeeper.SetUserRewardInfo(suite.ctx, mastercheftypes.UserRewardInfo{
					User:          positonAddress.String(),
					PoolId:        1,
					RewardDenom:   "uusdc",
					RewardPending: sdkmath.LegacyMustNewDecFromStr("100"),
				})
			},
		},
		{"positive case",
			&types.MsgClaimAllUserRewards{
				Sender: addresses[0].String(),
			},
			false,
			"",
			func() {
				suite.ResetSuite()
				suite.SetupCoinPrices(suite.ctx)
				initializeForClaimRewards(suite, addresses, asset1, asset2, true)
				openPosition(suite, addresses[0], collateralAmount, leverage)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
			_, err := msgServer.ClaimAllUserRewards(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
