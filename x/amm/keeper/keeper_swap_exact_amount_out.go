package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// InternalSwapExactAmountOut is a method for swapping to get an exact number of tokens out of a pool,
// using the provided swapFee.
// This is intended to allow different swap fees as determined by multi-hops,
// or when recovering from chain liveness failures.
func (k Keeper) InternalSwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	pool types.Pool,
	tokenInDenom string,
	tokenInMaxAmount math.Int,
	tokenOut sdk.Coin,
	swapFee math.LegacyDec,
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

	params := k.GetParams(ctx)
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	tokenIn, _, slippageAmount, weightBalanceBonus, oracleInAmount, swapFee, err := pool.SwapInAmtGivenOut(ctx, k.oracleKeeper, &snapshot, sdk.Coins{tokenOut}, tokenInDenom, swapFee, k.accountedPoolKeeper, math.LegacyOneDec(), params)
	if err != nil {
		return math.Int{}, err
	}
	tokenInAmount = tokenIn.Amount

	if tokenInAmount.LTE(math.ZeroInt()) {
		return math.Int{}, types.ErrTokenOutAmountZero
	}

	if tokenInAmount.GT(tokenInMaxAmount) {
		return math.Int{}, errorsmod.Wrapf(types.ErrLimitMaxAmount, "swap requires %s, which is greater than the amount %s", tokenIn, tokenInMaxAmount)
	}

	err = k.UpdatePoolForSwap(ctx, pool, sender, recipient, tokenIn, tokenOut, swapFee, oracleInAmount.ToInt(), math.ZeroInt(), weightBalanceBonus, true)
	if err != nil {
		return math.Int{}, err
	}

	// track slippage
	k.TrackSlippage(ctx, pool.PoolId, sdk.NewCoin(tokenIn.Denom, slippageAmount.ToInt()))

	return tokenInAmount, nil
}
