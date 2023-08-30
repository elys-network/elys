package keeper_test

import (
	"testing"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"

	margintypes "github.com/elys-network/elys/x/margin/types"
	"github.com/elys-network/elys/x/margin/types/mocks"
)

func TestCheckBalanceInvariant_InvalidBalance(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, ak, amm, oracle := app.MarginKeeper, app.AccountedPoolKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := sdk.NewCoins(sdk.NewCoin(ptypes.USDC, sdk.NewInt(100000)))
	// Mint 100000ATOM
	atomToken := sdk.NewCoins(sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000)))

	err := app.BankKeeper.MintCoins(ctx, types.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, types.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, addr[0], atomToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(100000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.USDC, sdk.NewInt(10000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.0")
	argExitFee := sdk.MustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
	}

	msg := types.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a ATOM+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(0))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a margin position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.USDC), sdk.NewInt(10000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(100000))

	// Create a margin position open msg
	msg2 := margintypes.NewMsgOpen(
		addr[0].String(),
		ptypes.USDC,
		sdk.NewInt(100),
		ptypes.ATOM,
		margintypes.Position_LONG,
		sdk.NewDec(5),
	)
	// Setup the mock checker
	mockAuthorizationChecker := new(mocks.AuthorizationChecker)
	mk.AuthorizationChecker = mockAuthorizationChecker

	// Mock behavior
	mockAuthorizationChecker.On("IsWhitelistingEnabled", ctx).Return(true)
	mockAuthorizationChecker.On("CheckIfWhitelisted", ctx, pool.GetAddress()).Return(true)

	// Setup the mock checker
	mockPositionChecker := new(mocks.PositionChecker)

	// Mock behavior
	mockPositionChecker.On("GetOpenMTPCount", ctx).Return(uint64(0))
	mockPositionChecker.On("GetMaxOpenPositions", ctx).Return(10)
	mk.PositionChecker = mockPositionChecker

	_, err = mk.Open(ctx, msg2)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.USDC), sdk.NewInt(10100))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(100000))

	// Check balance invariant check
	err = ak.InvariantCheck(ctx)
	require.NoError(t, err)

	mtpId := mtps[0].Id
	// Create a margin position close msg
	msg3 := margintypes.NewMsgClose(
		addr[0].String(),
		mtpId,
	)

	_, err = mk.Close(ctx, msg3)
	require.NoError(t, err)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.USDC), sdk.NewInt(10000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(100000))

	// Check balance invariant check
	err = ak.InvariantCheck(ctx)
	require.NoError(t, err)
}
