package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapExactAmountOut is a method for swapping to get an exact number of tokens out of a pool,
// using the provided swapFee.
// This is intended to allow different swap fees as determined by multi-hops,
// or when recovering from chain liveness failures.
func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	pool types.Pool,
	tokenInDenom string,
	tokenInMaxAmount math.Int,
	tokenOut sdk.Coin,
	swapFee sdk.Dec,
) (tokenInAmount math.Int, err error) {
	if tokenInDenom == tokenOut.Denom {
		return math.Int{}, errors.New("cannot trade same denomination in and out")
	}
	defer func() {
		if r := recover(); r != nil {
			tokenInAmount = math.Int{}
			err = fmt.Errorf("function swapExactAmountOut failed due to internal reason: %v", r)
		}
	}()

	poolOutBal := pool.GetTotalPoolLiquidity().AmountOf(tokenOut.Denom)
	if tokenOut.Amount.GTE(poolOutBal) {
		return math.Int{}, sdkerrors.Wrapf(types.ErrTooManyTokensOut,
			"can't get more tokens out than there are tokens in the pool")
	}

	tokenIn, err := pool.SwapInAmtGivenOut(ctx, sdk.Coins{tokenOut}, tokenInDenom, swapFee)
	if err != nil {
		return math.Int{}, err
	}
	tokenInAmount = tokenIn.Amount

	if tokenInAmount.LTE(sdk.ZeroInt()) {
		return math.Int{}, sdkerrors.Wrapf(types.ErrInvalidMathApprox, "token amount is zero or negative")
	}

	if tokenInAmount.GT(tokenInMaxAmount) {
		return math.Int{}, sdkerrors.Wrapf(types.ErrLimitMaxAmount, "Swap requires %s, which is greater than the amount %s", tokenIn, tokenInMaxAmount)
	}

	err = k.updatePoolForSwap(ctx, pool, sender, tokenIn, tokenOut)
	if err != nil {
		return math.Int{}, err
	}
	return tokenInAmount, nil
}
