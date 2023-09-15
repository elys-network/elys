package keeper_test

import (
	"fmt"
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	atypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCommitMintedLPTokenToCommitmentModule(t *testing.T) {
	app := simapp.InitElysTestApp(initChain)
	ctx := app.BaseApp.NewContext(initChain, tmproto.Header{})
	amm, bk := app.AmmKeeper, app.BankKeeper

	// Create Pool

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))
	transferAmt := sdk.NewCoin(ptypes.Elys, sdk.NewInt(100))

	// Deposit 100elys to FeeCollectorName wallet
	err := bk.SendCoinsFromAccountToModule(ctx, addr[0], authtypes.FeeCollectorName, sdk.NewCoins(transferAmt))
	require.NoError(t, err)

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000)))

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	var poolAssets []atypes.PoolAsset
	// Elys
	poolAssets = append(poolAssets, atypes.PoolAsset{
		Weight: sdk.NewInt(50),
		Token:  sdk.NewCoin(ptypes.Elys, sdk.NewInt(100000)),
	})

	// USDC
	poolAssets = append(poolAssets, atypes.PoolAsset{
		Weight: sdk.NewInt(50),
		Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(10000)),
	})

	argSwapFee, err := sdk.NewDecFromStr("0.1")
	require.NoError(t, err)

	argExitFee, err := sdk.NewDecFromStr("0.1")
	require.NoError(t, err)

	poolParams := &atypes.PoolParams{
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

	_, found := amm.GetPool(ctx, poolId)
	require.True(t, found)

	lpTokenDenom := types.GetPoolShareDenom(poolId)
	lpTokenBalance := bk.GetBalance(ctx, addr[0], lpTokenDenom)
	fmt.Println("lpTokenBalance")
	fmt.Println(lpTokenBalance)
}
