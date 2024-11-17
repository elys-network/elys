package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// MultihopSwapExactAmountOut defines the output denom and output amount for the last pool.
// Calculation starts by providing the tokenOutAmount of the final pool to calculate the required tokenInAmount
// the calculated tokenInAmount is used as defined tokenOutAmount of the previous pool, calculating in reverse order of the swap
// Transaction succeeds if the calculated tokenInAmount of the first pool is less than the defined tokenInMaxAmount defined.
func (k Keeper) RouteExactAmountOut(ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	routes []types.SwapAmountOutRoute,
	tokenInMaxAmount math.Int,
	tokenOut sdk.Coin,
	discount math.LegacyDec,
) (tokenInAmount math.Int, totalDiscountedSwapFee math.LegacyDec, discountOut math.LegacyDec, err error) {
	isMultiHopRouted, routeSwapFee, sumOfSwapFees := false, math.LegacyDec{}, math.LegacyDec{}
	route := types.SwapAmountOutRoutes(routes)
	if err := route.Validate(); err != nil {
		return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	defer func() {
		if r := recover(); r != nil {
			tokenInAmount = math.Int{}
			err = fmt.Errorf("function RouteExactAmountOut failed due to internal reason: %v", r)
		}
	}()

	// in this loop, we check if:
	// - the route is of length 2
	// - route 1 and route 2 don't trade via the same pool
	// - route 1 contains uosmo
	// - both route 1 and route 2 are incentivized pools
	// if all of the above is true, then we collect the additive and max fee between the two pools to later calculate the following:
	// total_swap_fee = total_swap_fee = max(swapfee1, swapfee2)
	// fee_per_pool = total_swap_fee * ((pool_fee) / (swapfee1 + swapfee2))
	if k.isElysRoutedMultihop(ctx, route, routes[0].TokenInDenom, tokenOut.Denom) {
		isMultiHopRouted = true
		routeSwapFee, sumOfSwapFees, err = k.getElysRoutedMultihopTotalSwapFee(ctx, route)
		if err != nil {
			return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee = math.LegacyZeroDec()

	// Determine what the estimated input would be for each pool along the multi-hop route
	// if we determined the route is an osmo multi-hop and both routes are incentivized,
	// we utilize a separate function that calculates the discounted swap fees
	var insExpected []math.Int
	if isMultiHopRouted {
		insExpected, err = k.createElysMultihopExpectedSwapOuts(ctx, routes, tokenOut, routeSwapFee, sumOfSwapFees)
	} else {
		insExpected, err = k.createMultihopExpectedSwapOuts(ctx, routes, tokenOut)
	}
	if err != nil {
		return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}
	if len(insExpected) == 0 {
		return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), nil
	}

	insExpected[0] = tokenInMaxAmount

	// Iterates through each routed pool and executes their respective swaps. Note that all of the work to get the return
	// value of this method is done when we calculate insExpected – this for loop primarily serves to execute the actual
	// swaps on each pool.
	for i, route := range routes {
		_tokenOut := tokenOut

		// If there is one pool left in the route, set the expected output of the current swap
		// to the estimated input of the final pool.
		if i != len(routes)-1 {
			_tokenOut = sdk.NewCoin(routes[i+1].TokenInDenom, insExpected[i+1])
		}

		// Execute the expected swap on the current routed pool
		pool, poolExists := k.GetPool(ctx, route.PoolId)
		if !poolExists {
			return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), types.ErrInvalidPoolId
		}

		// // check if pool is active, if not error
		// if !pool.IsActive(ctx) {
		// 	return math.Int{}, fmt.Errorf("pool %d is not active", pool.GetId())
		// }

		swapFee := pool.GetPoolParams().SwapFee
		if isMultiHopRouted && sumOfSwapFees.IsPositive() {
			swapFee = routeSwapFee.Mul((swapFee.Quo(sumOfSwapFees)))
		}

		// Apply discount to swap fee if applicable
		_, tier := k.tierKeeper.GetMembershipTier(ctx, sender)
		swapFee = types.ApplyDiscount(swapFee, tier.Discount)

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		_tokenInAmount, swapErr := k.InternalSwapExactAmountOut(ctx, sender, recipient, pool, route.TokenInDenom, insExpected[i], _tokenOut, swapFee)
		if swapErr != nil {
			return math.Int{}, math.LegacyZeroDec(), math.LegacyZeroDec(), swapErr
		}

		// Sets the final amount of tokens that need to be input into the first pool. Even though this is the final return value for the
		// whole method and will not change after the first iteration, we still iterate through the rest of the pools to execute their respective
		// swaps.
		if i == 0 {
			tokenInAmount = _tokenInAmount
		}
	}

	return tokenInAmount, totalDiscountedSwapFee, discount, nil
}
