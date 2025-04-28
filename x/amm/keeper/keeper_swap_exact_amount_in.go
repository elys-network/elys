package keeper

import (
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// InternalSwapExactAmountIn is an internal method for swapping an exact amount of tokens
// as input to a pool, using the provided swapFee. This is intended to allow
// different swap fees as determined by multi-hops, or when recovering from
// chain liveness failures.
// TODO: investigate if swapFee can be unexported
// https://github.com/osmosis-labs/osmosis/issues/3130
func (k Keeper) InternalSwapExactAmountIn(
	ctx sdk.Context,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	pool types.Pool,
	tokenIn sdk.Coin,
	tokenOutDenom string,
	tokenOutMinAmount math.Int,
	swapFee math.LegacyDec,
	takersFee math.LegacyDec,
) (tokenOutAmount math.Int, weightBalanceReward sdk.Coin, err error) {
	if tokenIn.Denom == tokenOutDenom {
		return math.Int{}, sdk.Coin{}, errors.New("cannot trade the same denomination in and out")
	}

	tokensIn := sdk.Coins{tokenIn}

	defer func() {
		if r := recover(); r != nil {
			tokenOutAmount = math.Int{}
			err = fmt.Errorf("function SwapExactAmountIn failed due to an internal reason: %v", r)
			ctx.Logger().Error(err.Error())
		}
	}()
	params := k.GetParams(ctx)
	// Executes the swap in the pool and stores the output. Updates pool assets but
	// does not actually transfer any tokens to or from the pool.
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	tokenOutCoin, _, slippageAmount, weightBalanceBonus, oracleOutAmount, swapFee, err := pool.SwapOutAmtGivenIn(ctx, k.oracleKeeper, &snapshot, tokensIn, tokenOutDenom, swapFee, k.accountedPoolKeeper, math.LegacyOneDec(), params, takersFee)
	if err != nil {
		return math.Int{}, sdk.Coin{}, err
	}

	tokenOutAmount = tokenOutCoin.Amount

	if !tokenOutAmount.IsPositive() {
		return math.Int{}, sdk.Coin{}, types.ErrTokenOutAmountZero
	}

	if tokenOutAmount.LT(tokenOutMinAmount) {
		return math.Int{}, sdk.Coin{}, errorsmod.Wrapf(types.ErrLimitMinAmount, "%s token is less than the minimum amount", tokenOutDenom)
	}

	// Settles balances between the tx sender and the pool to match the swap that was executed earlier.
	// Also emits a swap event and updates related liquidity metrics.
	bonusToken, err := k.UpdatePoolForSwap(ctx, pool, sender, recipient, tokenIn, tokenOutCoin, swapFee, slippageAmount, math.ZeroInt(), oracleOutAmount.TruncateInt(), weightBalanceBonus, takersFee, false)
	if err != nil {
		return math.Int{}, sdk.Coin{}, err
	}

	// track slippage
	k.TrackSlippage(ctx, pool.PoolId, sdk.NewCoin(tokenOutCoin.Denom, slippageAmount.RoundInt()))

	if pool.PoolParams.UseOracle {
		k.TrackWeightBreakingSlippage(ctx, pool.PoolId, sdk.NewCoin(tokenOutCoin.Denom, slippageAmount.RoundInt()))
	}

	return tokenOutAmount, bonusToken, nil
}
