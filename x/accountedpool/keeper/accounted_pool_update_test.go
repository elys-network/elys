package keeper_test

import (
	sdkmath "cosmossdk.io/math"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"

	"github.com/elys-network/elys/x/accountedpool/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestAccountedPoolUpdate(t *testing.T) {
	app := simapp.InitElysTestApp(true)
	ctx := app.BaseApp.NewContext(true)

	apk := app.AccountedPoolKeeper

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Initiate pool
	ammPool := ammtypes.Pool{
		PoolId:      0,
		Address:     addr[0].String(),
		PoolParams:  ammtypes.PoolParams{},
		TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
		PoolAssets: []ammtypes.PoolAsset{
			{Token: sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(100))},
			{Token: sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(1000))},
		},
		TotalWeight:       sdkmath.NewInt(100),
		RebalanceTreasury: addr[0].String(),
	}
	// Initiate pool
	accountedPool := types.AccountedPool{
		PoolId:      0,
		TotalShares: ammPool.TotalShares,
		PoolAssets:  []ammtypes.PoolAsset{},
		TotalWeight: ammPool.TotalWeight,
	}

	for _, asset := range ammPool.PoolAssets {
		accountedPool.PoolAssets = append(accountedPool.PoolAssets, asset)
	}
	// Set accounted pool
	apk.SetAccountedPool(ctx, accountedPool)

	perpetualPool := perpetualtypes.Pool{
		AmmPoolId:          0,
		Health:             sdkmath.LegacyNewDec(1),
		Enabled:            true,
		Closed:             false,
		BorrowInterestRate: sdkmath.LegacyNewDec(1),
		PoolAssetsLong: []perpetualtypes.PoolAsset{
			{
				Liabilities:         sdkmath.NewInt(400),
				Custody:             sdkmath.NewInt(0),
				AssetBalance:        sdkmath.NewInt(100),
				BlockBorrowInterest: sdkmath.NewInt(0),
				AssetDenom:          ptypes.BaseCurrency,
			},
			{
				Liabilities:         sdkmath.NewInt(0),
				Custody:             sdkmath.NewInt(50),
				AssetBalance:        sdkmath.NewInt(0),
				BlockBorrowInterest: sdkmath.NewInt(0),
				AssetDenom:          ptypes.ATOM,
			},
		},
		PoolAssetsShort: []perpetualtypes.PoolAsset{
			{
				Liabilities:         sdkmath.NewInt(400),
				Custody:             sdkmath.NewInt(0),
				AssetBalance:        sdkmath.NewInt(100),
				BlockBorrowInterest: sdkmath.NewInt(0),
				AssetDenom:          ptypes.BaseCurrency,
			},
			{
				Liabilities:         sdkmath.NewInt(0),
				Custody:             sdkmath.NewInt(50),
				AssetBalance:        sdkmath.NewInt(0),
				BlockBorrowInterest: sdkmath.NewInt(0),
				AssetDenom:          ptypes.ATOM,
			},
		},
	}
	// Update accounted pool
	apk.UpdateAccountedPool(ctx, ammPool, perpetualPool)

	apool, found := apk.GetAccountedPool(ctx, (uint64)(0))
	require.Equal(t, found, true)
	require.Equal(t, apool.PoolId, (uint64)(0))

	usdcBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.BaseCurrency)
	require.Equal(t, usdcBalance, sdkmath.NewInt(1000+400+100+400+100))
	atomBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.ATOM)
	require.Equal(t, atomBalance, sdkmath.NewInt(100))
}
