package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// RouteExactAmountIn defines the input denom and input amount for the first pool,
// the output of the first pool is chained as the input for the next routed pool
// transaction succeeds when final amount out is greater than tokenOutMinAmount defined.
func (k Keeper) RouteExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	routes []types.SwapAmountInRoute,
	tokenIn sdk.Coin,
	tokenOutMinAmount math.Int,
	discount sdk.Dec,
) (tokenOutAmount math.Int, err error) {
	var (
		isMultiHopRouted bool
		routeSwapFee     sdk.Dec
		sumOfSwapFees    sdk.Dec
	)

	route := types.SwapAmountInRoutes(routes)
	if err := route.Validate(); err != nil {
		return math.Int{}, err
	}

	// In this loop, we check if:
	// - the route is of length 2
	// - route 1 and route 2 don't trade via the same pool
	// - route 1 contains uelys
	// - both route 1 and route 2 are incentivized pools
	//
	// If all of the above is true, then we collect the additive and max fee between the
	// two pools to later calculate the following:
	// total_swap_fee = total_swap_fee = max(swapfee1, swapfee2)
	// fee_per_pool = total_swap_fee * ((pool_fee) / (swapfee1 + swapfee2))
	if k.isElysRoutedMultihop(ctx, route, routes[0].TokenOutDenom, tokenIn.Denom) {
		isMultiHopRouted = true
		routeSwapFee, sumOfSwapFees, err = k.getElysRoutedMultihopTotalSwapFee(ctx, route)
		if err != nil {
			return math.Int{}, err
		}
	}

	for i, route := range routes {
		// To prevent the multihop swap from being interrupted prematurely, we keep
		// the minimum expected output at a very low number until the last pool
		_outMinAmount := math.NewInt(1)
		if len(routes)-1 == i {
			_outMinAmount = tokenOutMinAmount
		}

		// Execute the expected swap on the current routed pool
		pool, poolExists := k.GetPool(ctx, route.PoolId)
		if !poolExists {
			return math.Int{}, types.ErrInvalidPoolId
		}

		// // check if pool is active, if not error
		// if !pool.IsActive(ctx) {
		// 	return math.Int{}, fmt.Errorf("pool %d is not active", pool.GetId())
		// }

		swapFee := pool.GetPoolParams().SwapFee

		// If we determined the route is an osmo multi-hop and both routes are incentivized,
		// we modify the swap fee accordingly.
		if isMultiHopRouted && sumOfSwapFees.IsPositive() {
			swapFee = routeSwapFee.Mul((swapFee.Quo(sumOfSwapFees)))
		}

		// Apply discount to swap fee if applicable
		swapFee, discount, err = k.ApplyDiscount(ctx, swapFee, discount, sender.String())
		if err != nil {
			return math.Int{}, err
		}

		tokenOutAmount, err = k.SwapExactAmountIn(ctx, sender, pool, tokenIn, route.TokenOutDenom, _outMinAmount, swapFee)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return math.Int{}, err
		}

		// Chain output of current pool as the input for the next routed pool
		tokenIn = sdk.NewCoin(route.TokenOutDenom, tokenOutAmount)
	}
	return tokenOutAmount, nil
}
