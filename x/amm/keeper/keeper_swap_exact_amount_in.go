package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// SwapExactAmountIn is an internal method for swapping an exact amount of tokens
// as input to a pool, using the provided swapFee. This is intended to allow
// different swap fees as determined by multi-hops, or when recovering from
// chain liveness failures.
// TODO: investigate if swapFee can be unexported
// https://github.com/osmosis-labs/osmosis/issues/3130
func (k Keeper) SwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	pool types.Pool,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	tokenOutMinAmount math.Int,
	swapFee sdk.Dec,
) (tokenOutAmount math.Int, err error) {
	if tokenIn.Denom == tokenOutDenom {
		return math.Int{}, errors.New("cannot trade the same denomination in and out")
	}

	tokensIn := sdk.Coins{tokenIn}

	defer func() {
		if r := recover(); r != nil {
			tokenOutAmount = math.Int{}
			err = fmt.Errorf("function SwapExactAmountIn failed due to an internal reason: %v", r)
		}
	}()

	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	snapshot := k.GetPoolSnapshotOrSet(ctx, pool)
	tokenOutCoin, _, slippageAmount, weightBalanceBonus, err := pool.SwapOutAmtGivenIn(ctx, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, swapFee, k.accountedPoolKeeper)
	if err != nil {
		return math.Int{}, err
	}

	tokenOutAmount = tokenOutCoin.Amount

	if !tokenOutAmount.IsPositive() {
		return math.Int{}, errorsmod.Wrapf(types.ErrInvalidMathApprox, "token amount must be positive")
	}

	if tokenOutAmount.LT(tokenOutMinAmount) {
		return math.Int{}, errorsmod.Wrapf(types.ErrLimitMinAmount, "%s token is less than the minimum amount", tokenOutDenom)
	}

	// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
	// Also emits a swap event and updates related liquidity metrics.
	swapOutFee, err := k.UpdatePoolForSwap(ctx, pool, sender, recipient, tokenIn, tokenOutCoin, sdk.ZeroDec(), swapFee, weightBalanceBonus)
	if err != nil {
		return math.Int{}, err
	}

	// track slippage
	k.TrackSlippage(ctx, pool.PoolId, sdk.NewCoin(tokenOutCoin.Denom, slippageAmount.RoundInt()))

	// Subtract swap out fee from the token out amount.
	return tokenOutAmount.Sub(swapOutFee), nil
}
