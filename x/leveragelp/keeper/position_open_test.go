package keeper_test

import (
	"cosmossdk.io/math"
	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/keeper"
	"github.com/elys-network/elys/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestOpenLong() {
	k := suite.app.LeveragelpKeeper
	suite.SetupCoinPrices(suite.ctx)
	addr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolAddr := ammtypes.NewPoolAddress(uint64(1))
	treasuryAddr := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	poolInit := sdk.Coins{sdk.NewInt64Coin("uusdc", 100000), sdk.NewInt64Coin("uusdt", 100000)}

	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, poolInit)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, poolAddr, poolInit)
	suite.Require().NoError(err)

	ammParams := suite.app.AmmKeeper.GetParams(suite.ctx)
	ammParams.WeightBreakingFeeMultiplier = math.LegacyZeroDec()
	ammParams.WeightBreakingFeeExponent = math.LegacyNewDecWithPrec(25, 1) // 2.5
	ammParams.WeightRecoveryFeePortion = math.LegacyNewDecWithPrec(10, 2)  // 10%
	ammParams.ThresholdWeightDifference = math.LegacyZeroDec()
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
	suite.Require().NoError(err)
	enablePoolMsg := types.MsgAddPool{
		Authority: authtypes.NewModuleAddress("gov").String(),
		Pool: types.AddPool{
			AmmPoolId:   1,
			LeverageMax: math.LegacyNewDec(10),
		},
	}
	msgServer := keeper.NewMsgServerImpl(*suite.app.LeveragelpKeeper)
	_, err = msgServer.AddPool(suite.ctx, &enablePoolMsg)
	suite.Require().NoError(err)
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
	})
	suite.Require().NoError(err)

	// open a position
	position, err := k.OpenLong(suite.ctx, &types.MsgOpen{
		Creator:          addr.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdkmath.NewInt(1000),
		AmmPoolId:        1,
		Leverage:         sdkmath.LegacyNewDec(5),
		StopLossPrice:    sdkmath.LegacyZeroDec(),
	})
	suite.Require().NoError(err)
	suite.Require().Equal(position.Address, addr.String())
	suite.Require().Equal(position.Collateral.String(), "1000uusdc")
	suite.Require().Equal(position.Liabilities.String(), "4000")
	// suite.Require().Equal(position.LeveragedLpAmount.String(), "49390000000000000") // slippage enabled on amm
	suite.Require().Equal(position.LeveragedLpAmount.String(), "50000000000000000") // slippage disabled on amm
	// suite.Require().Equal(position.PositionHealth.String(), "1.221000000000000000") // slippage enabled on amm
	suite.Require().Equal(position.PositionHealth.String(), "1.250000000000000000") // slippage disabled on amm
	suite.Require().Equal(position.Id, uint64(1))
	suite.Require().Equal(position.AmmPoolId, uint64(1))

	// add more to an existing position
	_, err = k.OpenConsolidate(suite.ctx, position, &types.MsgOpen{
		Creator:          addr.String(),
		CollateralAsset:  "uusdc",
		CollateralAmount: sdkmath.NewInt(1000),
		AmmPoolId:        1,
		Leverage:         sdkmath.LegacyNewDec(5),
		StopLossPrice:    sdkmath.LegacyZeroDec(),
	})
	suite.Require().NoError(err)
	position2, err := k.GetPosition(suite.ctx, position.GetOwnerAddress(), position.Id)
	suite.Require().NoError(err)
	suite.Require().Equal(position2.Address, addr.String())
	suite.Require().Equal(position2.Collateral.String(), "2000uusdc")
	suite.Require().Equal(position2.Liabilities.String(), "8000")
	// suite.Require().Equal(position2.LeveragedLpAmount.String(), "98805291560975610") // slippage enabled on amm
	suite.Require().Equal(position2.LeveragedLpAmount.String(), "100000000000000000") // slippage disabled on amm
	// suite.Require().Equal(position2.PositionHealth.String(), "1.210375000000000000") // slippage enabled on amm
	suite.Require().Equal(position2.PositionHealth.String(), "1.250000000000000000") // slippage disabled on amm
	suite.Require().Equal(position2.Id, uint64(1))
	suite.Require().Equal(position2.AmmPoolId, uint64(1))
}
