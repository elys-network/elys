package keeper_test

import (
	"time"

	sdkmath "cosmossdk.io/math"
	"github.com/cometbft/cometbft/crypto/ed25519"
	sdk "github.com/cosmos/cosmos-sdk/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	stablestakekeeper "github.com/elys-network/elys/v7/x/stablestake/keeper"
	stablestaketypes "github.com/elys-network/elys/v7/x/stablestake/types"
)

func (suite *KeeperTestSuite) TestQueryEstimation() {
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
				Token:                  poolInit[0],
				Weight:                 sdkmath.NewInt(10),
				ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
			},
			{
				Token:                  poolInit[1],
				Weight:                 sdkmath.NewInt(10),
				ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
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
	ammPool, found := suite.app.AmmKeeper.GetPool(suite.ctx, 1)
	suite.Require().True(found)
	err = suite.app.PerpetualKeeper.OnLeverageLpEnablePool(suite.ctx, ammPool)
	suite.Require().NoError(err)

	usdcToken := sdk.NewInt64Coin("uusdc", 100000)
	err = suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, sdk.Coins{usdcToken})
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, addr, sdk.Coins{usdcToken})
	suite.Require().NoError(err)

	leverageParams := suite.app.LeveragelpKeeper.GetParams(suite.ctx)
	leverageParams.EnabledPools = []uint64{1}
	err = suite.app.LeveragelpKeeper.SetParams(suite.ctx, &leverageParams)
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

	suite.AddBlockTime(time.Hour)

	estimation, err := k.CloseEst(suite.ctx, &types.QueryCloseEstRequest{
		Owner:    addr.String(),
		Id:       position.Id,
		LpAmount: sdkmath.NewInt(10000000000000000),
	})
	suite.Require().NoError(err)
	// Total liability is 4000, so 800 is the liability for 10000000000000000 Lp amount
	suite.Require().Equal(estimation.RepayAmount.String(), "809")
}
