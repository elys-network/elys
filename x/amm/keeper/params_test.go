package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	minttypes "github.com/cosmos/cosmos-sdk/x/mint/types"
	testkeeper "github.com/elys-network/elys/v6/testutil/keeper"
	"github.com/elys-network/elys/v6/x/amm/keeper"
	"github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestGetParams(t *testing.T) {
	k, ctx, _, _ := testkeeper.AmmKeeper(t)
	params := types.DefaultParams()
	params.BaseAssets = nil
	params.AllowedUpfrontSwapMakers = nil
	k.SetParams(ctx, params)

	require.EqualValues(t, params, k.GetParams(ctx))
}

func (suite *AmmKeeperTestSuite) TestCheckExistingPoolWithSameAssets() {
	suite.SetupTest()
	suite.SetupCoinPrices()
	suite.SetAmmParams()
	suite.SetupAssetProfile()

	// Bootstrap accounts
	sender := authtypes.NewModuleAddress(govtypes.ModuleName)
	params := suite.app.AmmKeeper.GetParams(suite.ctx)

	// Bootstrap balances
	poolCreationFee := sdk.NewCoin(ptypes.Elys, params.PoolCreationFee)
	coins := sdk.NewCoins(
		sdk.NewCoin("uatom", sdkmath.NewInt(1000)),
		sdk.NewCoin("uusdc", sdkmath.NewInt(1000)),
		poolCreationFee,
	)
	err := suite.app.BankKeeper.MintCoins(suite.ctx, minttypes.ModuleName, coins)
	suite.Require().NoError(err)
	err = suite.app.BankKeeper.SendCoinsFromModuleToAccount(suite.ctx, minttypes.ModuleName, sender, coins)
	suite.Require().NoError(err)

	// Create an initial pool
	msgServer := keeper.NewMsgServerImpl(*suite.app.AmmKeeper)
	poolAssets := []types.PoolAsset{
		{
			Token:                  sdk.NewCoin("uatom", sdkmath.NewInt(500)),
			Weight:                 sdkmath.NewInt(10),
			ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
		},
		{
			Token:                  sdk.NewCoin("uusdc", sdkmath.NewInt(500)),
			Weight:                 sdkmath.NewInt(10),
			ExternalLiquidityRatio: sdkmath.LegacyOneDec(),
		},
	}
	poolParams := types.PoolParams{
		SwapFee:   sdkmath.LegacyZeroDec(),
		UseOracle: false,
		FeeDenom:  ptypes.BaseCurrency,
	}
	_, err = msgServer.CreatePool(
		suite.ctx,
		&types.MsgCreatePool{
			Sender:     sender.String(),
			PoolParams: poolParams,
			PoolAssets: poolAssets,
		},
	)
	suite.Require().NoError(err)

	// Test case 1: No matching denoms
	newAssets := []types.PoolAsset{
		{
			Token: sdk.NewCoin("uosmo", sdkmath.NewInt(500)),
		},
		{
			Token: sdk.NewCoin("uusdc", sdkmath.NewInt(500)),
		},
	}
	exists := suite.app.AmmKeeper.CheckExistingPoolWithSameAssets(suite.ctx, newAssets)
	suite.Require().False(exists, "Expected no matching pool for new assets")

	// Test case 2: Matching denoms
	matchingAssets := []types.PoolAsset{
		{
			Token: sdk.NewCoin("uusdc", sdkmath.NewInt(500)),
		},
		{
			Token: sdk.NewCoin("uatom", sdkmath.NewInt(500)),
		},
	}
	exists = suite.app.AmmKeeper.CheckExistingPoolWithSameAssets(suite.ctx, matchingAssets)
	suite.Require().True(exists, "Expected a matching pool for the same assets")
}
