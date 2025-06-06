package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/amm/types"
)

// CalcOutRouteByDenom calculates the out route by denom
func (k Keeper) CalcOutRouteByDenom(ctx sdk.Context, denomOut string, denomIn string, baseCurrency string) ([]*types.SwapAmountOutRoute, error) {
	var route []*types.SwapAmountOutRoute

	// If the denoms are the same, throw an error
	if denomIn == denomOut {
		return nil, errorsmod.Wrap(types.ErrSameDenom, "denom in and denom out are the same")
	}

	// Check for a direct pool between the denoms
	if pool, found := k.GetBestPoolWithDenoms(ctx, []string{denomOut, denomIn}, false); found {
		// If the pool exists, return the route
		route = append(route, &types.SwapAmountOutRoute{
			PoolId:       pool.PoolId,
			TokenInDenom: denomIn,
		})
		return route, nil
	}

	// Find pool for initial denom to base currency
	pool, found := k.GetBestPoolWithDenoms(ctx, []string{denomOut, baseCurrency}, false)
	if !found {
		return nil, fmt.Errorf("no available pool for %s to base currency", denomOut)
	}
	// If the pool exists, append the route
	route = append(route, &types.SwapAmountOutRoute{
		PoolId:       pool.PoolId,
		TokenInDenom: baseCurrency,
	})

	// Find pool for base currency to target denom
	pool, found = k.GetBestPoolWithDenoms(ctx, []string{baseCurrency, denomIn}, false)
	if !found {
		return nil, fmt.Errorf("no available pool for base currency to %s", denomIn)
	}
	// If the pool exists, append the route
	route = append(route, &types.SwapAmountOutRoute{
		PoolId:       pool.PoolId,
		TokenInDenom: denomIn,
	})

	return route, nil
}
