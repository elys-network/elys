package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestCalcInRouteSpotPrice(t *testing.T) {
	k, ctx, accountedPoolKeeper, _ := keepertest.AmmKeeper(t)
	SetupMockPools(k, ctx)

	// Use token in for all tests
	tokenIn := sdk.NewCoin("denom1", sdk.NewInt(100))

	// Test single route
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom1").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom2").Return(sdk.NewInt(1000))
	routes := []*types.SwapAmountInRoute{{PoolId: 1, TokenOutDenom: "denom2"}}
	spotPrice, _, err := k.CalcInRouteSpotPrice(ctx, tokenIn, routes)
	require.NoError(t, err)
	require.NotZero(t, spotPrice)
	accountedPoolKeeper.AssertExpectations(t)

	// Test multiple routes
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(2), "denom1").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(2), "baseCurrency").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(3), "baseCurrency").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(3), "denom3").Return(sdk.NewInt(1000))
	routes = []*types.SwapAmountInRoute{
		{PoolId: 2, TokenOutDenom: "baseCurrency"},
		{PoolId: 3, TokenOutDenom: "denom3"},
	}
	spotPrice, _, err = k.CalcInRouteSpotPrice(ctx, tokenIn, routes)
	require.NoError(t, err)
	require.NotZero(t, spotPrice)
	accountedPoolKeeper.AssertExpectations(t)

	// Test no routes
	_, _, err = k.CalcInRouteSpotPrice(ctx, tokenIn, nil)
	require.Error(t, err)

	// Test invalid pool
	routes = []*types.SwapAmountInRoute{{PoolId: 9999, TokenOutDenom: "denom2"}}
	_, _, err = k.CalcInRouteSpotPrice(ctx, tokenIn, routes)
	require.Error(t, err)
}
