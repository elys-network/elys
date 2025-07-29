package keeper

import (
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// RouteExactAmountOut defines the output denom and output amount for the last pool.
// Calculation starts by providing the tokenOutAmount of the final pool to calculate the required tokenInAmount
// the calculated tokenInAmount is used as defined tokenOutAmount of the previous pool, calculating in reverse order of the swap
// Transaction succeeds if the calculated tokenInAmount of the first pool is less than the defined tokenInMaxAmount defined.
func (k Keeper) RouteExactAmountOut(ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	routes []types.SwapAmountOutRoute,
	tokenInMaxAmount math.Int,
	tokenOut sdk.Coin,
) (tokenInAmount math.Int, totalDiscountedSwapFee osmomath.BigDec, discountOut osmomath.BigDec, err error) {
	route := types.SwapAmountOutRoutes(routes)
	if err := route.Validate(); err != nil {
		return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	defer func() {
		if r := recover(); r != nil {
			tokenInAmount = math.Int{}
			err = fmt.Errorf("function RouteExactAmountOut failed due to internal reason: %v", r)
			ctx.Logger().Error(err.Error())
		}
	}()

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee = osmomath.ZeroBigDec()

	insExpected, err := k.createMultihopExpectedSwapOuts(ctx, routes, tokenOut)
	if err != nil {
		return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}
	if len(insExpected) == 0 {
		return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), nil
	}

	insExpected[0] = tokenInMaxAmount

	_, tier := k.tierKeeper.GetMembershipTier(ctx, sender)
	discount := tier.Discount

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
			return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), types.ErrInvalidPoolId
		}

		// // check if pool is active, if not error
		// if !pool.IsActive(ctx) {
		// 	return math.Int{}, fmt.Errorf("pool %d is not active", pool.GetId())
		// }

		swapFee := pool.GetPoolParams().GetBigDecSwapFee().Quo(osmomath.NewBigDec(int64(len(routes))))
		takersFee := k.parameterKeeper.GetParams(ctx).GetBigDecTakerFees().Quo(osmomath.NewBigDec(int64(len(routes))))

		// Apply discount to swap fee if applicable
		swapFee = types.ApplyDiscount(swapFee, osmomath.BigDecFromDec(discount))

		// Calculate the total discounted swap fee
		totalDiscountedSwapFee = totalDiscountedSwapFee.Add(swapFee)

		_tokenInAmount, weightBalanceReward, swapErr := k.InternalSwapExactAmountOut(ctx, sender, recipient, pool, route.TokenInDenom, insExpected[i], _tokenOut, swapFee, takersFee)
		if swapErr != nil {
			return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), swapErr
		}

		if weightBalanceReward.Amount.IsPositive() && weightBalanceReward.Denom == route.TokenInDenom {
			// If the weight balance reward is positive, we need to add it to the tokenInAmount
			tokenInAmount = tokenInAmount.Add(weightBalanceReward.Amount)
		}

		// Sets the final amount of tokens that need to be input into the first pool. Even though this is the final return value for the
		// whole method and will not change after the first iteration, we still iterate through the rest of the pools to execute their respective
		// swaps.
		if i == 0 {
			tokenInAmount = _tokenInAmount
		}
	}

	return tokenInAmount, totalDiscountedSwapFee, osmomath.BigDecFromDec(discount), nil
}
