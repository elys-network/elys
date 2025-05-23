package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// RouteExactAmountIn defines the input denom and input amount for the first pool,
// the output of the first pool is chained as the input for the next routed pool
// transaction succeeds when final amount out is greater than tokenOutMinAmount defined.
func (k Keeper) RouteExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	routes []types.SwapAmountInRoute,
	tokenIn sdk.Coin,
	tokenOutMinAmount math.Int,
) (tokenOutAmount math.Int, totalDiscountedSwapFee osmomath.BigDec, discountOut osmomath.BigDec, err error) {
	route := types.SwapAmountInRoutes(routes)
	if err := route.Validate(); err != nil {
		return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
	}

	// Initialize the total discounted swap fee
	totalDiscountedSwapFee = osmomath.ZeroBigDec()

	_, tier := k.tierKeeper.GetMembershipTier(ctx, sender)
	discount := tier.Discount
	var weightBalanceReward sdk.Coin
	for i, route := range routes {
		// recipient is the same as the sender until the last pool
		actualRecipient := sender
		if len(routes)-1 == i {
			actualRecipient = recipient
		}

		// To prevent the multihop swap from being interrupted prematurely, we keep
		// the minimum expected output at a very low number until the last pool
		_outMinAmount := math.NewInt(1)
		if len(routes)-1 == i {
			_outMinAmount = tokenOutMinAmount
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

		tokenOutAmount, weightBalanceReward, err = k.InternalSwapExactAmountIn(ctx, sender, actualRecipient, pool, tokenIn, route.TokenOutDenom, _outMinAmount, swapFee, takersFee)
		if err != nil {
			ctx.Logger().Error(err.Error())
			return math.Int{}, osmomath.ZeroBigDec(), osmomath.ZeroBigDec(), err
		}

		if weightBalanceReward.Amount.IsPositive() && weightBalanceReward.Denom == route.TokenOutDenom {
			// If the weight balance reward is positive, we need to add it to the tokenOutAmount
			tokenOutAmount = tokenOutAmount.Add(weightBalanceReward.Amount)
		}

		// Chain output of current pool as the input for the next routed pool
		tokenIn = sdk.NewCoin(route.TokenOutDenom, tokenOutAmount)
	}

	return tokenOutAmount, totalDiscountedSwapFee, osmomath.BigDecFromDec(tier.Discount), nil
}
