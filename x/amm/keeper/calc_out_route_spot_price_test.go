package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	keepertest "github.com/elys-network/elys/testutil/keeper"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/stretchr/testify/require"
)

func TestCalcOutRouteSpotPrice(t *testing.T) {
	k, ctx, accountedPoolKeeper, _ := keepertest.AmmKeeper(t)
	SetupMockPools(k, ctx)

	// Test single route
	tokenOut := sdk.NewCoin("denom2", sdk.NewInt(100))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom2").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(1), "denom1").Return(sdk.NewInt(1000))
	routes := []*types.SwapAmountOutRoute{{PoolId: 1, TokenInDenom: "denom1"}}
	spotPrice, _, _, _, _, err := k.CalcOutRouteSpotPrice(ctx, tokenOut, routes, sdk.ZeroDec())
	require.NoError(t, err)
	require.NotZero(t, spotPrice)
	accountedPoolKeeper.AssertExpectations(t)

	// Test multiple routes
	tokenOut = sdk.NewCoin("denom3", sdk.NewInt(100))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(3), "denom3").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(3), "baseCurrency").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(2), "baseCurrency").Return(sdk.NewInt(1000))
	accountedPoolKeeper.On("GetAccountedBalance", ctx, uint64(2), "denom1").Return(sdk.NewInt(1000))
	routes = []*types.SwapAmountOutRoute{
		{PoolId: 3, TokenInDenom: "baseCurrency"},
		{PoolId: 2, TokenInDenom: "denom1"},
	}
	spotPrice, _, _, _, _, err = k.CalcOutRouteSpotPrice(ctx, tokenOut, routes, sdk.ZeroDec())
	require.NoError(t, err)
	require.NotZero(t, spotPrice)
	accountedPoolKeeper.AssertExpectations(t)

	// Test no routes
	_, _, _, _, _, err = k.CalcOutRouteSpotPrice(ctx, tokenOut, nil, sdk.ZeroDec())
	require.Error(t, err)

	// Test invalid pool
	routes = []*types.SwapAmountOutRoute{{PoolId: 9999, TokenInDenom: "denom2"}}
	_, _, _, _, _, err = k.CalcOutRouteSpotPrice(ctx, tokenOut, routes, sdk.ZeroDec())
	require.Error(t, err)
}
