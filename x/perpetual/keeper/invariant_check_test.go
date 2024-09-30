package keeper_test

import (
	"cosmossdk.io/math"
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/amm/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"

	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
)

func TestCheckBalanceInvariant_InvalidBalance(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, math.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, math.NewInt(100000000000))}

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
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, math.NewInt(1000000000)),
		},
		{
			Weight: math.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(10000000000)),
		},
	}

	argSwapFee := math.LegacyMustNewDecFromStr("0.0")
	argExitFee := math.LegacyMustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		SwapFee:   argSwapFee,
		ExitFee:   argExitFee,
		UseOracle: true,
	}

	msg := types.NewMsgCreatePool(
		addr[0].String(),
		poolParams,
		poolAssets,
	)

	// Create a ATOM+USDC pool
	poolId, err := amm.CreatePool(ctx, msg)
	require.NoError(t, err)
	require.Equal(t, poolId, uint64(1))

	pools := amm.GetAllPool(ctx)

	// check length of pools
	require.Equal(t, len(pools), 1)

	// check block height
	require.Equal(t, int64(0), ctx.BlockHeight())

	pool, found := amm.GetPool(ctx, poolId)
	require.Equal(t, found, true)

	poolAddress := sdk.MustAccAddressFromBech32(pool.GetAddress())
	require.NoError(t, err)

	// Balance check before create a perpetual position
	balances := app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(10000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(1000000000))

	// Create a perpetual position open msg
	msg2 := perpetualtypes.NewMsgOpen(
		addr[0].String(),
		perpetualtypes.Position_LONG,
		math.LegacyNewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, math.NewInt(100000000)),
		math.LegacyMustNewDecFromStr(perpetualtypes.TakeProfitPriceDefault),
		math.LegacyNewDec(100),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(10100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(1000000000))

	// Check balance invariant check
	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	mtpId := mtps[0].Id
	// Create a perpetual position close msg
	msg3 := perpetualtypes.NewMsgClose(
		addr[0].String(),
		mtpId,
		mtps[0].Custody,
	)

	_, err = mk.Close(ctx, msg3)
	require.NoError(t, err)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), math.NewInt(10100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), math.NewInt(1000000000))

	// Check balance invariant check
	err = mk.InvariantCheck(ctx)
	// TODO:
	// Need to fix invariant balance check function
	require.Equal(t, err, errors.New("balance mismatch!"))
}
