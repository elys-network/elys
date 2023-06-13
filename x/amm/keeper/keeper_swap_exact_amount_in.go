package keeper

import (
	"errors"
	"fmt"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
)

// swapExactAmountIn is an internal method for swapping an exact amount of tokens
// as input to a pool, using the provided swapFee. This is intended to allow
// different swap fees as determined by multi-hops, or when recovering from
// chain liveness failures.
// TODO: investigate if swapFee can be unexported
// https://github.com/osmosis-labs/osmosis/issues/3130
func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	pool types.Pool,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	tokenOutMinAmount math.Int,
) (tokenOutAmount math.Int, err error) {
	if tokenIn.Denom == tokenOutDenom {
		return math.Int{}, errors.New("cannot trade same denomination in and out")
	}
	tokensIn := sdk.Coins{tokenIn}

	defer func() {
		if r := recover(); r != nil {
			tokenOutAmount = math.Int{}
			err = fmt.Errorf("function swapExactAmountIn failed due to internal reason: %v", r)
		}
	}()

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	tokenOutCoin, weightBalanceBonus, err := pool.SwapOutAmtGivenIn(ctx, k.oracleKeeper, tokensIn, tokenOutDenom)
	if err != nil {
		return math.Int{}, err
	}

	tokenOutAmount = tokenOutCoin.Amount

	if !tokenOutAmount.IsPositive() {
		return math.Int{}, sdkerrors.Wrapf(types.ErrInvalidMathApprox, "token amount must be positive")
	}

	if tokenOutAmount.LT(tokenOutMinAmount) {
		return math.Int{}, sdkerrors.Wrapf(types.ErrLimitMinAmount, "%s token is lesser than min amount", tokenOutDenom)
	}

	// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
	// Also emits swap event and updates related liquidity metrics
	if err := k.UpdatePoolForSwap(ctx, pool, sender, tokenIn, tokenOutCoin, weightBalanceBonus); err != nil {
		return math.Int{}, err
	}

	return tokenOutAmount, nil
}
