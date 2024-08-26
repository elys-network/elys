package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	mastercheftypes "github.com/elys-network/elys/x/masterchef/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func initializeForClaimRewards(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string, createAmmPool bool) {
	fee := sdk.MustNewDecFromStr("0.0002")
	issueAmount := sdk.NewInt(10_000_000_000_000)
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewCoin(ptypes.ATOM, issueAmount),
			sdk.NewCoin(ptypes.Elys, issueAmount),
			sdk.NewCoin(ptypes.BaseCurrency, issueAmount),
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
		PoolParams: &ammtypes.PoolParams{
			SwapFee:                     fee,
			ExitFee:                     fee,
			UseOracle:                   true,
			WeightBreakingFeeMultiplier: fee,
			WeightBreakingFeeExponent:   fee,
			ExternalLiquidityRatio:      fee,
			WeightRecoveryFeePortion:    fee,
			ThresholdWeightDifference:   fee,
			FeeDenom:                    ptypes.Elys,
		},
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  sdk.NewInt64Coin(asset1, 100_000_000),
				Weight: sdk.NewInt(50),
			},
			{
				Token:  sdk.NewInt64Coin(asset2, 100_000_000),
				Weight: sdk.NewInt(50),
			},
		},
	}
	if createAmmPool {
		poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
		if err != nil {
			panic(err)
		}

		suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(poolId))

	}
	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  issueAmount.QuoRaw(20),
	}
	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err := stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
	msgBond.Creator = addresses[2].String()
	_, err = stableStakeMsgServer.Bond(suite.ctx, &msgBond)
	if err != nil {
		panic(err)
	}
}

func (suite *KeeperTestSuite) TestMsgServerClaimRewards() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	leverage := sdk.MustNewDecFromStr("2.0")
	collateralAmount := sdk.NewInt(10000000)
	testCases := []struct {
		name                 string
		input                *types.MsgClaimRewards
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"position not found",
			&types.MsgClaimRewards{
				Sender: addresses[0].String(),
				Ids:    []uint64{1},
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
				Sender: addresses[0].String(),
				Ids:    []uint64{1},
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
					RewardPending: sdk.MustNewDecFromStr("100"),
				})
			},
		},
		{"positive case",
			&types.MsgClaimRewards{
				Sender: addresses[0].String(),
				Ids:    []uint64{1},
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
			msgServer := keeper.NewMsgServerImpl(suite.app.LeveragelpKeeper)
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
