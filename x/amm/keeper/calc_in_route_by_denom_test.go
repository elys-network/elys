package keeper_test

import (
	"testing"

	simapp "github.com/elys-network/elys/v4/app"
	"github.com/stretchr/testify/require"
)

func TestCalcInRouteByDenom(t *testing.T) {
	app := simapp.InitElysTestApp(initChain, t)
	ctx := app.BaseApp.NewContext(initChain)
	k := app.AmmKeeper

	// Setup mock pools and assets
	SetupMockPools(k, ctx)

	// Test direct pool route
	route, err := k.CalcInRouteByDenom(ctx, "denom1", "denom2", "uusdc")
	require.NoError(t, err)
	require.Len(t, route, 1)
	require.Equal(t, route[0].TokenOutDenom, "denom2")

	// Test route via base currency
	route, err = k.CalcInRouteByDenom(ctx, "denom1", "denom3", "uusdc")
	require.NoError(t, err)
	require.Len(t, route, 2)
	require.Equal(t, route[0].TokenOutDenom, "uusdc")
	require.Equal(t, route[1].TokenOutDenom, "denom3")

	// Test no available pool
	_, err = k.CalcInRouteByDenom(ctx, "denom1", "nonexistent", "uusdc")
	require.Error(t, err)

	// Test same input and output denomination
	route, err = k.CalcInRouteByDenom(ctx, "denom1", "denom1", "uusdc")
	require.NoError(t, err)
	require.Len(t, route, 1)
	require.Equal(t, route[0].TokenOutDenom, "denom1")
}
