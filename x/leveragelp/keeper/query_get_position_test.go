package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestQueryGetPosition() {
	k := suite.app.LeveragelpKeeper
	suite.SetupCoinPrices(suite.ctx)
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolAddr := ammtypes.NewPoolAddress(uint64(1))
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	params := stablestaketypes.DefaultParams()
	suite.app.StablestakeKeeper.SetParams(suite.ctx, params)

	pool := types.Pool{
		AmmPoolId:         1,
		Health:            sdkmath.LegacyZeroDec(),
		LeveragedLpAmount: sdkmath.ZeroInt(),
		LeverageMax:       sdkmath.LegacyOneDec().MulInt64(10),
	}
	poolInit := sdk.Coins{sdk.NewInt64Coin("uusdc", 100000), sdk.NewInt64Coin("uusdt", 100000)}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInit)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInit)
	suite.Require().NoError(err)

	ammParams := suite.app.AmmKeeper.GetParams(suite.ctx)
	ammParams.WeightBreakingFeeMultiplier = sdkmath.LegacyZeroDec()
	ammParams.WeightBreakingFeeExponent = sdkmath.LegacyNewDecWithPrec(25, 1) // 2.5
	ammParams.WeightRecoveryFeePortion = sdkmath.LegacyNewDecWithPrec(10, 2)  // 10%
	ammParams.ThresholdWeightDifference = sdkmath.LegacyZeroDec()
	suite.app.AmmKeeper.SetParams(suite.ctx, ammParams)

	suite.app.AmmKeeper.SetPool(suite.ctx, ammtypes.Pool{
		PoolId:            1,
		Address:           poolAddr.String(),
		RebalanceTreasury: treasuryAddr.String(),
		PoolParams: ammtypes.PoolParams{
			SwapFee:   sdkmath.LegacyZeroDec(),
			UseOracle: true,
			FeeDenom:  "uusdc",
		},
		TotalShares: sdk.NewCoin("amm/pool/1", sdkmath.NewInt(2).Mul(ammtypes.OneShare)),
		PoolAssets: []ammtypes.PoolAsset{
			{
				Token:  poolInit[0],
				Weight: sdkmath.NewInt(10),
			},
			{
				Token:  poolInit[1],
				Weight: sdkmath.NewInt(10),
			},
		},
		TotalWeight: sdkmath.NewInt(20),
	})
	k.SetPool(suite.ctx, pool)
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdc",
		Liquidity: sdkmath.NewInt(100000),
	})
	suite.app.AmmKeeper.SetDenomLiquidity(suite.ctx, ammtypes.DenomLiquidity{
		Denom:     "uusdt",
		Liquidity: sdkmath.NewInt(100000),
	})

	usdcToken := sdk.NewInt64Coin("uusdc", 100000)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	stableMsgServer := stablestakekeeper.NewMsgServerImpl(*suite.app.StablestakeKeeper)
	_, err = stableMsgServer.Bond(suite.ctx, &stablestaketypes.MsgBond{
		Creator: addr.String(),
		Amount:  sdkmath.NewInt(10000),
		PoolId:  1,
	})
	suite.Require().NoError(err)

	// open a position
	position, _ := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdkmath.NewInt(1000),
		AmmPoolId:        1,
		Leverage:         sdkmath.LegacyNewDec(5),
		StopLossPrice:    sdkmath.LegacyZeroDec(),
	}, 1)

	res, _ := k.Position(suite.ctx, &types.PositionRequest{Address: addr.String(), Id: position.Id})
	updated_leverage := sdkmath.LegacyNewDec(5)

	suite.Require().Equal(position, res.Position.Position)
	suite.Require().Equal(updated_leverage, res.Position.UpdatedLeverage)

	expected := types.PositionAndInterest{
		Position: &types.QueryPosition{
			Position:         position,
			UpdatedLeverage:  updated_leverage,
			PositionUsdValue: sdkmath.LegacyNewDec(5000).Quo(sdkmath.LegacyNewDec(1000000)),
		},
		InterestRateHour:    sdkmath.LegacyMustNewDecFromStr("0.000017123287671233"),
		InterestRateHourUsd: sdkmath.LegacyZeroDec(),
	}
	pos_for_address_res, _ := k.QueryPositionsForAddress(suite.ctx, &types.PositionsForAddressRequest{Address: addr.String(), Pagination: nil})

	suite.Require().Equal(expected.Position, pos_for_address_res.Positions[0].Position)
	suite.Require().True(expected.InterestRateHour.Equal(pos_for_address_res.Positions[0].InterestRateHour))
	suite.Require().True(expected.InterestRateHourUsd.Equal(pos_for_address_res.Positions[0].InterestRateHourUsd))
}
