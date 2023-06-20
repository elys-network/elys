package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
)

func TestCollectGasFeesToIncentiveModule(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})

	ik, bk, amm := app.IncentiveKeeper, app.BankKeeper, app.AmmKeeper
	// Collect gas fees
	collectedAmt := ik.CollectGasFeesToIncentiveModule(ctx)

	// rewards should be zero
	require.True(t, collectedAmt.IsZero())

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))
	transferAmt := sdk.NewCoin(ptypes.Elys, sdk.NewInt(100))

	// Deposit 100elys to FeeCollectorName wallet
	err := bk.SendCoinsFromAccountToModule(ctx, addr[0], authtypes.FeeCollectorName, sdk.NewCoins(transferAmt))
	require.NoError(t, err)

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.USDC, sdk.NewInt(100000)))

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	var poolAssets []ammtypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdk.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdk.NewInt(100000)),
	})

	// USDC
	poolAssets = append(poolAssets, ammtypes.PoolAsset{
		Weight: sdk.NewInt(50),
		Token:  sdk.NewCoin(ptypes.USDC, sdk.NewInt(10000)),
	})

	argSwapFee, err := sdk.NewDecFromStr("0.1")
	require.NoError(t, err)

	argExitFee, err := sdk.NewDecFromStr("0.1")
	require.NoError(t, err)

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
	}

	msg := types.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a Elys+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(0))
	//

	// Collect gas fees again
	collectedAmt = ik.CollectGasFeesToIncentiveModule(ctx)

	// It should be 10 usdc
	require.Equal(t, collectedAmt, sdk.Coins{sdk.NewCoin(ptypes.USDC, sdk.NewInt(10))})
}
