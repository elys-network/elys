package keeper_test

import (
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/v6/app"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"

	"github.com/elys-network/elys/v6/x/accountedpool/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/v6/x/perpetual/types"
	"github.com/stretchr/testify/require"
)

func TestAccountedPoolUpdate(t *testing.T) {
	app := simapp.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	apk := app.AccountedPoolKeeper

	err := simapp.SetStakingParam(app, ctx)
	require.NoError(t, err)

	// Generate 1 random account with 1000stake balanced
	addr := simapp.AddTestAddrs(app, ctx, 1, sdkmath.NewInt(1000000))

	// Initiate pool
	ammPool := ammtypes.Pool{
		PoolId:      0,
		Address:     addr[0].String(),
		PoolParams:  ammtypes.PoolParams{},
		TotalShares: sdk.NewCoin("lp-token", sdkmath.NewInt(100)),
		PoolAssets: []ammtypes.PoolAsset{
			{Token: sdk.NewCoin(ptypes.ATOM, sdkmath.NewInt(5000))},
			{Token: sdk.NewCoin(ptypes.BaseCurrency, sdkmath.NewInt(1000))},
		},
		TotalWeight:       sdkmath.NewInt(100),
		RebalanceTreasury: addr[0].String(),
	}
	// Initiate pool
	accountedPool := types.AccountedPool{
		PoolId:           0,
		TotalTokens:      []sdk.Coin{},
		NonAmmPoolTokens: []sdk.Coin{},
	}

	for _, asset := range ammPool.PoolAssets {
		accountedPool.TotalTokens = append(accountedPool.TotalTokens, asset.Token)
	}

	// Set accounted pool
	apk.SetAccountedPool(ctx, accountedPool)

	perpetualPool := perpetualtypes.Pool{
		AmmPoolId:                  0,
		BaseAssetLiabilitiesRatio:  sdkmath.LegacyZeroDec(),
		QuoteAssetLiabilitiesRatio: sdkmath.LegacyZeroDec(),
		BorrowInterestRate:         sdkmath.LegacyNewDec(1),
		PoolAssetsLong: []perpetualtypes.PoolAsset{
			{
				Liabilities: sdkmath.NewInt(400),
				Custody:     sdkmath.NewInt(50),
				AssetDenom:  ptypes.BaseCurrency,
			},
			{
				Liabilities: sdkmath.NewInt(0),
				Custody:     sdkmath.NewInt(50),
				AssetDenom:  ptypes.ATOM,
			},
		},
		PoolAssetsShort: []perpetualtypes.PoolAsset{
			{
				Liabilities: sdkmath.NewInt(400),
				Custody:     sdkmath.NewInt(70),
				AssetDenom:  ptypes.BaseCurrency,
			},
			{
				Liabilities: sdkmath.NewInt(0),
				Custody:     sdkmath.NewInt(50),
				AssetDenom:  ptypes.ATOM,
			},
		},
	}
	// Update accounted pool
	err = apk.PerpetualUpdates(ctx, ammPool, perpetualPool)
	require.NoError(t, err)

	apool, found := apk.GetAccountedPool(ctx, (uint64)(0))
	require.Equal(t, found, true)
	require.Equal(t, apool.PoolId, (uint64)(0))

	usdcBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.BaseCurrency)
	require.Equal(t, usdcBalance, sdkmath.NewInt(1000+400-50+400-70))
	atomBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.ATOM)
	require.Equal(t, atomBalance, sdkmath.NewInt(5000-50-50))
}
