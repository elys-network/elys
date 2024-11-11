package keeper_test

import (
	"cosmossdk.io/math"
	"testing"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	simapp "github.com/elys-network/elys/app"
	ammtypes "github.com/elys-network/elys/x/amm/types"

	"github.com/elys-network/elys/x/accountedpool/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
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
		TotalShares:      ammPool.TotalShares,
		PoolAssets:       []ammtypes.PoolAsset{},
		TotalWeight:      ammPool.TotalWeight,
		NonAmmPoolTokens: []sdk.Coin{},
	}

	for _, asset := range ammPool.PoolAssets {
		accountedPool.PoolAssets = append(accountedPool.PoolAssets, asset)
	}
	// Set accounted pool
	apk.SetAccountedPool(ctx, accountedPool)

	perpetualPool := perpetualtypes.Pool{
		AmmPoolId:          0,
		Health:             sdkmath.LegacyNewDec(1),
		BorrowInterestRate: sdkmath.LegacyNewDec(1),
		PoolAssetsLong: []perpetualtypes.PoolAsset{
			{
				Liabilities:           math.NewInt(400),
				Custody:               math.NewInt(50),
				TakeProfitCustody:     math.NewInt(10),
				TakeProfitLiabilities: math.NewInt(20),
				AssetDenom:            ptypes.BaseCurrency,
			},
			{
				Liabilities:           math.NewInt(0),
				Custody:               math.NewInt(50),
				TakeProfitCustody:     math.ZeroInt(),
				TakeProfitLiabilities: math.ZeroInt(),
				AssetDenom:            ptypes.ATOM,
			},
		},
		PoolAssetsShort: []perpetualtypes.PoolAsset{
			{
				Liabilities:           math.NewInt(400),
				Custody:               math.NewInt(70),
				TakeProfitCustody:     math.ZeroInt(),
				TakeProfitLiabilities: math.ZeroInt(),
				AssetDenom:            ptypes.BaseCurrency,
			},
			{
				Liabilities:           math.NewInt(0),
				Custody:               math.NewInt(50),
				TakeProfitCustody:     math.ZeroInt(),
				TakeProfitLiabilities: math.ZeroInt(),
				AssetDenom:            ptypes.ATOM,
			},
		},
	}
	// Update accounted pool
	err = apk.PerpetualUpdates(ctx, ammPool, perpetualPool, false)
	require.NoError(t, err)

	apool, found := apk.GetAccountedPool(ctx, (uint64)(0))
	require.Equal(t, found, true)
	require.Equal(t, apool.PoolId, (uint64)(0))

	usdcBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.BaseCurrency)
	require.Equal(t, usdcBalance, sdkmath.NewInt(1000+400-50+400-70))
	atomBalance := apk.GetAccountedBalance(ctx, (uint64)(0), ptypes.ATOM)
	require.Equal(t, atomBalance, sdkmath.NewInt(5000-50-50))
}
