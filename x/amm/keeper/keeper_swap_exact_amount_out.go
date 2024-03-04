package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapExactAmountOut is a method for swapping to get an exact number of tokens out of a pool,
// using the provided swapFee.
// This is intended to allow different swap fees as determined by multi-hops,
// or when recovering from chain liveness failures.
func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	pool types.Pool,
	tokenInDenom string,
	tokenInMaxAmount math.Int,
	tokenOut sdk.Coin,
	swapFee sdk.Dec,
) (tokenInAmount math.Int, err error) {
	if tokenInDenom == tokenOut.Denom {
		return math.Int{}, errors.New("cannot trade the same denomination in and out")
	}

	defer func() {
		if r := recover(); r != nil {
			tokenInAmount = math.Int{}
			err = fmt.Errorf("function SwapExactAmountOut failed due to an internal reason: %v", r)
		}
	}()

	poolOutBal := pool.GetTotalPoolLiquidity().AmountOf(tokenOut.Denom)
	if tokenOut.Amount.GTE(poolOutBal) {
		return math.Int{}, errorsmod.Wrapf(types.ErrTooManyTokensOut, "cannot get more tokens out than there are tokens in the pool")
	}

	snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
	tokenIn, _, slippageAmount, weightBalanceBonus, err := pool.SwapInAmtGivenOut(ctx, k.oracleKeeper, &snapshot, sdk.Coins{tokenOut}, tokenInDenom, swapFee, k.accountedPoolKeeper)
	if err != nil {
		return math.Int{}, err
	}
	tokenInAmount = tokenIn.Amount

	if tokenInAmount.LTE(sdk.ZeroInt()) {
		return math.Int{}, errorsmod.Wrapf(types.ErrInvalidMathApprox, "token amount is zero or negative")
	}

	if tokenInAmount.GT(tokenInMaxAmount) {
		return math.Int{}, errorsmod.Wrapf(types.ErrLimitMaxAmount, "swap requires %s, which is greater than the amount %s", tokenIn, tokenInMaxAmount)
	}

	_, err = k.UpdatePoolForSwap(ctx, pool, sender, recipient, tokenIn, tokenOut, swapFee, sdk.ZeroDec(), weightBalanceBonus)
	if err != nil {
		return math.Int{}, err
	}

	// track slippage
	k.TrackSlippage(ctx, pool.PoolId, sdk.NewCoin(tokenIn.Denom, slippageAmount.RoundInt()))
	return tokenInAmount, nil
}
