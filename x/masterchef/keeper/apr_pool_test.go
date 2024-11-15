package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammkeeper "github.com/elys-network/elys/x/amm/keeper"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCalculatePoolAprs(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	simapp.SetMasterChefParams(app, ctx)
	err := simapp.SetStakingParam(app, ctx)
	require.NoError(t, err)
	simapp.SetupAssetProfile(app, ctx)

	mk, amm, oracle := app.MasterchefKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(100010))

	// Create a pool
	// Mint 100000USDC + 10 ELYS (pool creation fee)
	coins := sdk.NewCoins(sdk.NewInt64Coin(ptypes.Elys, 10000000), sdk.NewInt64Coin(ptypes.BaseCurrency, 100000))
	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, coins)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], coins)
	require.NoError(t, err)

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(1000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100)),
	})

	poolParams := &ammtypes.PoolParams{
		SwapFee:                     sdkmath.LegacyZeroDec(),
		ExitFee:                     sdkmath.LegacyZeroDec(),
		UseOracle:                   false,
		WeightBreakingFeeMultiplier: sdkmath.LegacyZeroDec(),
		WeightBreakingFeeExponent:   sdkmath.LegacyNewDecWithPrec(25, 1), // 2.5
		WeightRecoveryFeePortion:    sdkmath.LegacyNewDecWithPrec(10, 2), // 10%
		ThresholdWeightDifference:   sdkmath.LegacyZeroDec(),
		FeeDenom:                    "",
	}

	// Create a Elys+USDC pool
	msgServer := ammkeeper.NewMsgServerImpl(*amm)
	resp, err := msgServer.CreatePool(
		ctx,
		&ammtypes.MsgCreatePool{
			Sender:     addr[0].String(),
			PoolParams: poolParams,
			PoolAssets: poolAssets,
		})

	require.NoError(t, err)
	require.Equal(t, resp.PoolID, uint64(1))

	poolInfo, found := mk.GetPoolInfo(ctx, resp.PoolID)
	require.True(t, found)

	poolInfo.DexApr = sdkmath.LegacyNewDecWithPrec(1, 2)  // 1%
	poolInfo.EdenApr = sdkmath.LegacyNewDecWithPrec(2, 2) // 2%
	mk.SetPoolInfo(ctx, poolInfo)

	// When passing empty array
	aprs := mk.CalculatePoolAprs(ctx, []uint64{})
	require.Len(t, aprs, 2) // setting it 2 because PoolId = math.MaxInt16 gets initiated in EndBlock
	require.Equal(t, aprs[0].TotalApr.String(), "0.030000000000000000")

	// When passing specific id
	aprs = mk.CalculatePoolAprs(ctx, []uint64{1})
	require.Len(t, aprs, 1)
	require.Equal(t, aprs[0].TotalApr.String(), "0.030000000000000000")

	// When passing invalid id
	aprs = mk.CalculatePoolAprs(ctx, []uint64{4})
	require.Len(t, aprs, 1)
	require.Equal(t, aprs[0].TotalApr.String(), "0.000000000000000000")
}
