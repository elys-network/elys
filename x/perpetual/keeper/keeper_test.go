package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"

	"github.com/cometbft/cometbft/crypto/ed25519"
	oraclekeeper "github.com/elys-network/elys/x/oracle/keeper"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func TestSetGetMTP(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	for i := 0; i < 2; i++ {
		mtp := types.MTP{
			Address:                        addr[i].String(),
			CollateralAsset:                ptypes.BaseCurrency,
			CustodyAsset:                   "ATOM",
			Collateral:                     math.NewInt(0),
			Liabilities:                    math.NewInt(0),
			BorrowInterestPaidCollateral:   math.NewInt(0),
			BorrowInterestPaidCustody:      math.NewInt(0),
			BorrowInterestUnpaidCollateral: math.NewInt(0),
			Custody:                        math.NewInt(0),
			MtpHealth:                      math.LegacyNewDec(0),
			Position:                       types.Position_LONG,
			Id:                             0,
			ConsolidateLeverage:            math.LegacyZeroDec(),
			SumCollateral:                  math.ZeroInt(),
		}
		err := perpetual.SetMTP(ctx, &mtp)
		require.NoError(t, err)
	}

	mtpCount := perpetual.GetMTPCount(ctx)
	require.Equal(t, mtpCount, (uint64)(2))
}

func TestGetAllWhitelistedAddress(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	perpetual := app.PerpetualKeeper

	// Generate 2 random accounts with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 2, math.NewInt(1000000))

	// Set whitelisted addresses
	perpetual.WhitelistAddress(ctx, addr[0])
	perpetual.WhitelistAddress(ctx, addr[1])

	// Get all whitelisted addresses
	whitelisted := perpetual.GetAllWhitelistedAddress(ctx)

	// length should be 2
	require.Equal(t, len(whitelisted), 2)

	// If addr[0] is whitelisted
	require.Contains(t,
		whitelisted,
		addr[0],
	)

	// If addr[1] is whitelisted
	require.Contains(t,
		whitelisted,
		addr[1],
	)
}

func SetupStableCoinPrices(ctx sdk.Context, oracle oraclekeeper.Keeper) {
	// prices set for USDT and USDC
	provider := sdk.AccAddress(ed25519.GenPrivKey().PubKey().Address())
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.BaseCurrency,
		Display: "USDC",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   "uusdt",
		Display: "USDT",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.Elys,
		Display: "ELYS",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   ptypes.ATOM,
		Display: "ATOM",
		Decimal: 6,
	})
	oracle.SetAssetInfo(ctx, oracletypes.AssetInfo{
		Denom:   "uatom",
		Display: "uatom",
		Decimal: 6,
	})

	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDC",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "USDT",
		Price:     math.LegacyNewDec(1),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ELYS",
		Price:     math.LegacyNewDec(23),
		Source:    "elys",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "ATOM",
		Price:     math.LegacyNewDec(6),
		Source:    "atom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})
	oracle.SetPrice(ctx, oracletypes.Price{
		Asset:     "uatom",
		Price:     math.LegacyNewDec(1),
		Source:    "uatom",
		Provider:  provider.String(),
		Timestamp: uint64(ctx.BlockTime().Unix()),
	})

}
