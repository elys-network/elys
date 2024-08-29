package keeper_test

import (
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (suite KeeperTestSuite) TestQueryGetPosition() {
	k := suite.app.LeveragelpKeeper
	suite.SetupCoinPrices(suite.ctx)
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	params := stablestaketypes.DefaultParams()
	suite.app.StablestakeKeeper.SetParams(suite.ctx, params)

	pool := types.Pool{
		AmmPoolId:         1,
		Enabled:           true,
		Closed:            false,
		Health:            sdk.ZeroDec(),
		LeveragedLpAmount: sdk.ZeroInt(),
		LeverageMax:       sdk.OneDec().MulInt64(10),
	}
	poolInit := sdk.Coins{sdk.NewInt64Coin("uusdc", 100000), sdk.NewInt64Coin("uusdt", 100000)}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInit)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInit)
	suite.Require().NoError(err)

	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:                     sdk.ZeroDec(),
			ExitFee:                     sdk.ZeroDec(),
			UseOracle:                   true,
			WeightBreakingFeeMultiplier: sdk.ZeroDec(),
			WeightBreakingFeeExponent:   sdk.NewDecWithPrec(25, 1), // 2.5
			ExternalLiquidityRatio:      sdk.NewDec(1),
			WeightRecoveryFeePortion:    sdk.NewDecWithPrec(10, 2), // 10%
			ThresholdWeightDifference:   sdk.ZeroDec(),
			FeeDenom:                    "uusdc",
		},
		TotalShares: sdk.NewCoin("amm/pool/1", sdk.NewInt(2).Mul(ammtypes.OneShare)),
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  poolInit[0],
				Weight: sdk.NewInt(10),
			},
			{
				Token:  poolInit[1],
				Weight: sdk.NewInt(10),
			},
		},
		TotalWeight: sdk.NewInt(20),
	})
	k.SetPool(suite.ctx, pool)
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdc",
		Liquidity: sdk.NewInt(100000),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdt",
		Liquidity: sdk.NewInt(100000),
	})

	usdcToken := sdk.NewInt64Coin("uusdc", 100000)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	stableMsgServer := stablestakekeeper.NewMsgServerImpl(suite.app.StablestakeKeeper)
	_, err = stableMsgServer.Bond(sdk.WrapSDKContext(suite.ctx), &stablestaketypes.MsgBond{
		Creator: addr.String(),
		Amount:  sdk.NewInt(10000),
	})
	suite.Require().NoError(err)

	// open a position
	position, _ := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdk.NewInt(1000),
		AmmPoolId:        1,
		Leverage:         sdk.NewDec(5),
		StopLossPrice:    sdk.ZeroDec(),
	})

	res, _ := k.Position(suite.ctx, &types.PositionRequest{Address: addr.String(), Id: position.Id})
	updated_leverage := sdk.NewDec(5)

	suite.Require().Equal(position, res.Position.Position)
	suite.Require().Equal(updated_leverage, res.Position.UpdatedLeverage)


	expected := types.PositionAndInterest{
		Position: &types.QueryPosition{
			Position: position,
			UpdatedLeverage: updated_leverage,
		},
		InterestRateHour: sdk.MustNewDecFromStr("0.000017123287671233"),
		InterestRateHourUsd: sdk.ZeroDec(),
	}
	pos_for_address_res, _ := k.QueryPositionsForAddress(suite.ctx, &types.PositionsForAddressRequest{Address: addr.String(), Pagination: nil} )

	suite.Require().Equal(expected.Position, pos_for_address_res.Positions[0].Position)
	suite.Require().Equal(expected.InterestRateHour, pos_for_address_res.Positions[0].InterestRateHour)
}
