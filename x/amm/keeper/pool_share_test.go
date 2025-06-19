package keeper_test

import (
	"fmt"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/v6/app"
	"github.com/elys-network/elys/v6/x/amm/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCommitMintedLPTokenToCommitmentModule(t *testing.T) {
	app := simapp.InitElysTestApp(initChain, t)
	ctx := app.BaseApp.NewContext(initChain)
	amm, bk := app.AmmKeeper, app.BankKeeper

	simapp.SetupAssetProfile(app, ctx)
	// Create Pool
	err := simapp.SetStakingParam(app, ctx)
	require.NoError(t, err)
	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))
	transferAmt := sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100))

	// Deposit 100elys to FeeCollectorName wallet
	err = bk.SendCoinsFromAccountToModule(ctx, addr[0], authtypes.FeeCollectorName, sdk.NewCoins(transferAmt))
	require.NoError(t, err)

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(100000)))

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	var poolAssets []types.PoolAsset
	// Elys
	poolAssets = append(poolAssets, types.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdkmath.NewInt(100000)),
	})

	// USDC
	poolAssets = append(poolAssets, types.PoolAsset{
		Weight: sdkmath.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(10000)),
	})

	argSwapFee, err := sdkmath.LegacyNewDecFromStr("0.01")
	require.NoError(t, err)

	poolParams := types.PoolParams{
		SwapFee: argSwapFee,
	}

	msg := types.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a Elys+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	_, found := amm.GetPool(ctx, poolId)
	require.True(t, found)

	lpTokenDenom := types.GetPoolShareDenom(poolId)
	lpTokenBalance := bk.GetBalance(ctx, addr[0], lpTokenDenom)
	fmt.Println("lpTokenBalance", lpTokenBalance.String())
}
