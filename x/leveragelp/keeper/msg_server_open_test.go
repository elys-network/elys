package keeper_test

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
)

func initialize(suite *KeeperTestSuite, addresses []sdk.AccAddress) {
	fee := sdk.MustNewDecFromStr("0.0002")
	SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
	for _, address := range addresses {
		coins := sdk.NewCoins(
			sdk.NewInt64Coin(ptypes.ATOM, 10000000000000),
			sdk.NewInt64Coin(ptypes.Elys, 10000000000000),
			sdk.NewInt64Coin(ptypes.BaseCurrency, 10000000000000),
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
				Token:  sdk.NewInt64Coin(ptypes.ATOM, 10000000/2),
				Weight: sdk.NewInt(50),
			},
			{
				Token:  sdk.NewInt64Coin(ptypes.BaseCurrency, 10000000/2),
				Weight: sdk.NewInt(50),
			},
		},
	}
	poolId, err := suite.app.AmmKeeper.CreatePool(suite.ctx, &msgCreatePool)
	if err != nil {
		panic(err)
	}
	suite.app.LeveragelpKeeper.SetPool(suite.ctx, types.NewPool(poolId))
	msgBond := stabletypes.MsgBond{
		Creator: addresses[1].String(),
		Amount:  sdk.NewInt(10000000000000 / 2),
	}
	stableStakeMsgServer := stablekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
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

func (suite *KeeperTestSuite) TestOpen() {
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	initialize(suite, addresses)
	testCases := []struct {
		name                 string
		input                *types.MsgOpen
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"failed user authorization",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  "stake",
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			errorsmod.Wrap(types.ErrUnauthorised, "unauthorised").Error(),
			func() {
				suite.EnableWhiteListing()
			},
		},
		{"no positions allowed",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  "stake",
				CollateralAmount: sdk.NewInt(1000000000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("10.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			true,
			errorsmod.Wrap(types.ErrMaxOpenPositions, "cannot open new positions").Error(),
			func() {
				suite.DisableWhiteListing()
				suite.SetMaxOpenPositions(0)
			},
		},
		{"Open Long",
			&types.MsgOpen{
				Creator:          addresses[0].String(),
				CollateralAsset:  ptypes.BaseCurrency,
				CollateralAmount: sdk.NewInt(1000),
				AmmPoolId:        1,
				Leverage:         sdk.MustNewDecFromStr("2.0"),
				StopLossPrice:    sdk.MustNewDecFromStr("100.0"),
			},
			false,
			"",
			func() {
				suite.SetMaxOpenPositions(5)
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
