package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/amm/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestCalcInRouteByDenom(t *testing.T) {
	k, ctx := keepertest.AmmKeeper(t)

	// Setup mock pools and assets
	setupMockPools(k, ctx)

	// Test direct pool route
	route, err := k.CalcInRouteByDenom(ctx, "denom1", "denom2", "baseCurrency")
	require.NoError(t, err)
	require.Len(t, route, 1)
	require.Equal(t, route[0].TokenOutDenom, "denom2")

	// Test route via base currency
	route, err = k.CalcInRouteByDenom(ctx, "denom1", "denom3", "baseCurrency")
	require.NoError(t, err)
	require.Len(t, route, 2)
	require.Equal(t, route[0].TokenOutDenom, "baseCurrency")
	require.Equal(t, route[1].TokenOutDenom, "denom3")

	// Test no available pool
	_, err = k.CalcInRouteByDenom(ctx, "denom1", "nonexistent", "baseCurrency")
	require.Error(t, err)

	// Test same input and output denomination
	_, err = k.CalcInRouteByDenom(ctx, "denom1", "denom1", "baseCurrency")
	require.Error(t, err)
}

func setupMockPools(k *keeper.Keeper, ctx sdk.Context) {
	// Create and set mock pools
	pools := []types.Pool{
		{
			PoolId: 1,
			PoolAssets: []types.PoolAsset{
				{Token: sdk.Coin{Denom: "denom1", Amount: sdk.NewInt(1000)}},
				{Token: sdk.Coin{Denom: "denom2", Amount: sdk.NewInt(1000)}},
			},
		},
		{
			PoolId: 2,
			PoolAssets: []types.PoolAsset{
				{Token: sdk.Coin{Denom: "denom1", Amount: sdk.NewInt(1000)}},
				{Token: sdk.Coin{Denom: "baseCurrency", Amount: sdk.NewInt(1000)}},
			},
		},
		{
			PoolId: 3,
			PoolAssets: []types.PoolAsset{
				{Token: sdk.Coin{Denom: "baseCurrency", Amount: sdk.NewInt(1000)}},
				{Token: sdk.Coin{Denom: "denom3", Amount: sdk.NewInt(1000)}},
			},
		},
	}

	for _, pool := range pools {
		k.SetPool(ctx, pool)
	}
}
