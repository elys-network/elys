package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/perpetual/types"

	tmproto "github.com/cometbft/cometbft/proto/tendermint/types"
	simapp "github.com/elys-network/elys/app"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	"github.com/stretchr/testify/require"
)

func TestCheckAndLiquidateUnhealthyPosition(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true, tmproto.Header{})

	mk, amm, oracle := app.PerpetualKeeper, app.AmmKeeper, app.OracleKeeper

	// Setup coin prices
	SetupStableCoinPrices(ctx, oracle)

	// Set asset profile
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.BaseCurrency,
		Denom:     ptypes.BaseCurrency,
		Decimals:  6,
	})
	app.AssetprofileKeeper.SetEntry(ctx, assetprofiletypes.Entry{
		BaseDenom: ptypes.ATOM,
		Denom:     ptypes.ATOM,
		Decimals:  6,
	})

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdk.NewInt(1000000000000))

	// Create a pool
	// Mint 100000USDC
	usdcToken := []sdk.Coin{sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(200000000000))}
	// Mint 100000ATOM
	atomToken := []sdk.Coin{sdk.NewCoin(ptypes.ATOM, sdk.NewInt(200000000000))}

	err := app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, usdcToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], usdcToken)
	require.NoError(t, err)

	err = app.BankKeeper.MintCoins(ctx, ammtypes.ModuleName, atomToken)
	require.NoError(t, err)
	err = app.BankKeeper.SendCoinsFromModuleToAccount(ctx, ammtypes.ModuleName, addr[0], atomToken)
	require.NoError(t, err)

	poolAssets := []ammtypes.PoolAsset{
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.ATOM, sdk.NewInt(10000000000)),
		},
		{
			Weight: sdk.NewInt(50),
			Token:  sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000000)),
		},
	}

	argSwapFee := sdk.MustNewDecFromStr("0.0")
	argExitFee := sdk.MustNewDecFromStr("0.0")

	poolParams := &ammtypes.PoolParams{
		SwapFee: argSwapFee,
		ExitFee: argExitFee,
	}

	msg := ammtypes.NewMsgCreatePool(
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
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100000000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))

	// Create a perpetual position open msg
	msg2 := types.NewMsgOpen(
		addr[0].String(),
		types.Position_LONG,
		sdk.NewDec(5),
		ptypes.ATOM,
		sdk.NewCoin(ptypes.BaseCurrency, sdk.NewInt(100000000)),
		sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
	)

	_, err = mk.Open(ctx, msg2, false)
	require.NoError(t, err)

	mtps := mk.GetAllMTPs(ctx)
	require.Equal(t, len(mtps), 1)

	balances = app.BankKeeper.GetAllBalances(ctx, poolAddress)
	require.Equal(t, balances.AmountOf(ptypes.BaseCurrency), sdk.NewInt(100100000000))
	require.Equal(t, balances.AmountOf(ptypes.ATOM), sdk.NewInt(10000000000))

	_, found = mk.OpenLongChecker.GetPool(ctx, pool.PoolId)
	require.Equal(t, found, true)

	err = mk.InvariantCheck(ctx)
	require.Equal(t, err, nil)

	mtp := mtps[0]

	perpPool, _ := mk.GetPool(ctx, pool.PoolId)
	err = mk.CheckAndLiquidateUnhealthyPosition(ctx, &mtp, perpPool, pool, ptypes.BaseCurrency, 6)
	require.NoError(t, err)

	// Set borrow interest rate to 100% to test liquidation
	perpPool.BorrowInterestRate = sdk.MustNewDecFromStr("1.0")
	mk.SetPool(ctx, perpPool)

	// Check MTP
	require.Equal(t, types.MTP{
		Address:                        addr[0].String(),
		CollateralAsset:                "uusdc",
		TradingAsset:                   "uatom",
		LiabilitiesAsset:               "uusdc",
		CustodyAsset:                   "uatom",
		Collateral:                     sdk.NewInt(100000000),
		Liabilities:                    sdk.NewInt(400000000),
		BorrowInterestPaidCollateral:   sdk.NewInt(0),
		BorrowInterestPaidCustody:      sdk.NewInt(0),
		BorrowInterestUnpaidCollateral: sdk.NewInt(0),
		Custody:                        sdk.NewInt(49751243),
		TakeProfitLiabilities:          sdk.NewInt(495049497),
		TakeProfitCustody:              sdk.NewInt(49751243),
		Leverage:                       sdk.NewDec(5),
		MtpHealth:                      sdk.MustNewDecFromStr("1.249999982500000000"),
		Position:                       types.Position_LONG,
		Id:                             uint64(1),
		AmmPoolId:                      uint64(1),
		ConsolidateLeverage:            sdk.NewDec(4),
		SumCollateral:                  sdk.NewInt(100000000),
		TakeProfitPrice:                sdk.MustNewDecFromStr(types.TakeProfitPriceDefault),
		TakeProfitBorrowRate:           sdk.MustNewDecFromStr("1.0"),
		FundingFeePaidCollateral:       sdk.NewInt(0),
		FundingFeePaidCustody:          sdk.NewInt(0),
		FundingFeeReceivedCollateral:   sdk.NewInt(0),
		FundingFeeReceivedCustody:      sdk.NewInt(0),
		OpenPrice:                      sdk.MustNewDecFromStr("10.050000157785002477"),
	}, mtp)
}
