package keeper_test

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	stablekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stabletypes "github.com/elys-network/elys/x/stablestake/types"
	"time"
)

func initializeForClose(suite *KeeperTestSuite, addresses []sdk.AccAddress, asset1, asset2 string) {
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
				Token:  sdk.NewInt64Coin(asset2, 1000_000_000),
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
		Amount:  issueAmount.QuoRaw(20),
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

func (suite *KeeperTestSuite) TestClose() {
	suite.ResetSuite()
	SetupCoinPrices(suite.ctx, suite.app.OracleKeeper)
	addresses := simapp.AddTestAddrs(suite.app, suite.ctx, 10, sdk.NewInt(1000000))
	asset1 := ptypes.ATOM
	asset2 := ptypes.BaseCurrency
	initializeForClose(suite, addresses, asset1, asset2)
	testCases := []struct {
		name                 string
		input                *types.MsgClose
		expectErr            bool
		expectErrMsg         string
		prerequisiteFunction func()
	}{
		{"No position to close",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: sdk.NewInt(0),
			},
			true,
			types.ErrPositionDoesNotExist.Error(),
			func() {
			},
		},
		{"Unlock time not reached",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: sdk.NewInt(0),
			},
			true,
			"insufficient withdrawable tokens",
			func() {
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: sdk.NewInt(10000000),
					AmmPoolId:        1,
					Leverage:         sdk.MustNewDecFromStr("2.0"),
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				if err != nil {
					panic(err)
				}
			},
		},
		{"Closing whole position",
			&types.MsgClose{
				Creator:  addresses[0].String(),
				Id:       1,
				LpAmount: sdk.MustNewDecFromStr("9997999787743811730").TruncateInt(),
			},
			false,
			"",
			func() {
				poolB, _ := suite.app.LeveragelpKeeper.GetAmmPool(suite.ctx, 1)
				fmt.Println("Before Open total shares: " + poolB.TotalShares.String())
				tvlB, _ := poolB.TVL(suite.ctx, suite.app.OracleKeeper)
				fmt.Println("Before Open TVL: " + tvlB.String())
				msg := types.MsgOpen{
					Creator:          addresses[0].String(),
					CollateralAsset:  ptypes.BaseCurrency,
					CollateralAmount: sdk.NewInt(10000000),
					AmmPoolId:        1,
					Leverage:         sdk.MustNewDecFromStr("2.0"),
					StopLossPrice:    sdk.MustNewDecFromStr("50.0"),
				}
				_, err := suite.app.LeveragelpKeeper.Open(suite.ctx, &msg)
				poolA, _ := suite.app.LeveragelpKeeper.GetAmmPool(suite.ctx, 1)
				fmt.Println("After Open total shares: " + poolA.TotalShares.String())
				tvlA, _ := poolA.TVL(suite.ctx, suite.app.OracleKeeper)
				fmt.Println("After Open TVL: " + tvlA.String())
				if err != nil {
					panic(err)
				}
				position, err := suite.app.LeveragelpKeeper.GetPosition(suite.ctx, addresses[0].String(), 1)
				if err != nil {
					panic(err)
				}
				fmt.Println("LP AMOUNT: " + position.LeveragedLpAmount.String())
				fmt.Println("LP AMOUNT VALUE: " + tvlA.Mul(position.LeveragedLpAmount.ToLegacyDec()).Quo(poolA.TotalShares.Amount.ToLegacyDec()).String())
				suite.AddBlockTime(time.Hour)
				fmt.Println("Sender: " + addresses[0].String())
			},
		},
	}

	for _, tc := range testCases {
		suite.Run(tc.name, func() {
			tc.prerequisiteFunction()
			_, err := suite.app.LeveragelpKeeper.Close(suite.ctx, tc.input)
			if tc.expectErr {
				suite.Require().Error(err)
				suite.Require().Contains(err.Error(), tc.expectErrMsg)
			} else {
				suite.Require().NoError(err)
			}
		})
	}
}
