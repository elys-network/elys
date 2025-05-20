package keeper_test

import (
	"testing"

	simapp "github.com/elys-network/elys/v4/app"
	"github.com/stretchr/testify/require"
)

func TestCalcOutRouteByDenom(t *testing.T) {
	t.Parallel()
	app := simapp.InitElysTestApp(initChain, t)
	ctx := app.BaseApp.NewContext(initChain)
	k := app.AmmKeeper

	// Setup mock pools and assets
	SetupMockPools(k, ctx)

	// Test direct pool route
	route, err := k.CalcOutRouteByDenom(ctx, "denom2", "denom1", "uusdc")
	require.NoError(t, err)
	require.Len(t, route, 1)
	require.Equal(t, route[0].TokenInDenom, "denom1")

	// Test route via base currency
	route, err = k.CalcOutRouteByDenom(ctx, "denom3", "denom1", "uusdc")
	require.NoError(t, err)
	require.Len(t, route, 2)
	require.Equal(t, route[0].TokenInDenom, "uusdc")
	require.Equal(t, route[1].TokenInDenom, "denom1")

	// Test no available pool
	_, err = k.CalcOutRouteByDenom(ctx, "nonexistent", "denom1", "uusdc")
	require.Error(t, err)

	// Test same input and output denomination
	_, err = k.CalcOutRouteByDenom(ctx, "denom1", "denom1", "uusdc")
	require.Error(t, err)
}
