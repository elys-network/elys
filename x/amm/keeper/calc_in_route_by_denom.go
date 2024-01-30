package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// CalcInRouteByDenom calculates the in route by denom
func (k Keeper) CalcInRouteByDenom(ctx sdk.Context, denomIn string, denomOut string, baseCurrency string) ([]*types.SwapAmountInRoute, error) {
	var route []*types.SwapAmountInRoute

	// Check for a direct pool between the denoms
	if pool, found := k.GetBestPoolWithDenoms(ctx, []string{denomIn, denomOut}); found {
		// If the pool exists, return the route
		route = append(route, &types.SwapAmountInRoute{
			PoolId:        pool.PoolId,
			TokenOutDenom: denomOut,
		})
		return route, nil
	}

	// Find pool for initial denom to base currency
	pool, found := k.GetBestPoolWithDenoms(ctx, []string{denomIn, baseCurrency})
	if !found {
		return nil, fmt.Errorf("no available pool for %s to base currency", denomIn)
	}
	// If the pool exists, append the route
	route = append(route, &types.SwapAmountInRoute{
		PoolId:        pool.PoolId,
		TokenOutDenom: baseCurrency,
	})

	// Find pool for base currency to target denom
	pool, found = k.GetBestPoolWithDenoms(ctx, []string{baseCurrency, denomOut})
	if !found {
		return nil, fmt.Errorf("no available pool for base currency to %s", denomOut)
	}
	// If the pool exists, append the route
	route = append(route, &types.SwapAmountInRoute{
		PoolId:        pool.PoolId,
		TokenOutDenom: denomOut,
	})

	return route, nil
}
