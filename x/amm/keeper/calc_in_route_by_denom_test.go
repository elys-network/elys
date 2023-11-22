package keeper_test

import (
	"testing"

	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/stretchr/testify/require"
)

func TestCalcInRouteByDenom(t *testing.T) {
	k, ctx, _, _ := keepertest.AmmKeeper(t)

	// Setup mock pools and assets
	SetupMockPools(k, ctx)

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
