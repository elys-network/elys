package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// JoinPoolNoSwap aims to LP exactly enough to pool #{poolId} to get shareOutAmount number of LP shares.
// If the required tokens is greater than tokenInMaxs, returns an error & the message reverts.
// Leftover tokens that weren't LP'd (due to being at inexact ratios) remain in the sender account.
//
// JoinPoolNoSwap determines the maximum amount that can be LP'd without any swap,
// by looking at the ratio of the total LP'd assets. (e.g. 2 osmo : 1 atom)
// It then finds the maximal amount that can be LP'd.
func (k Keeper) JoinPoolNoSwap(
	ctx sdk.Context,
	sender sdk.AccAddress,
	poolId uint64,
	shareOutAmount sdkmath.Int,
	tokenInMaxs sdk.Coins,
) (tokenIn sdk.Coins, sharesOut sdkmath.Int, err error) {
	// defer to catch panics, in case something internal overflows.
	defer func() {
		if r := recover(); r != nil {
			tokenIn = sdk.Coins{}
			sharesOut = sdkmath.Int{}
			err = fmt.Errorf("function JoinPoolNoSwap failed due to internal reason: %v", r)
		}
	}()
	// all pools handled within this method are pointer references, `JoinPool` directly updates the pools
	pool, poolExists := k.GetPool(ctx, poolId)
	if !poolExists {
		return nil, sdkmath.ZeroInt(), types.ErrInvalidPoolId
	}

	if !pool.PoolParams.UseOracle {
		tokensIn := tokenInMaxs
		if len(tokensIn) != 1 {
			// we do an abstract calculation on the lp liquidity coins needed to have
			// the designated amount of given shares of the pool without performing swap
			neededLpLiquidity, err := pool.GetMaximalNoSwapLPAmount(shareOutAmount)
			if err != nil {
				return nil, sdkmath.ZeroInt(), err
			}

			// check that needed lp liquidity does not exceed the given `tokenInMaxs` parameter. Return error if so.
			// if tokenInMaxs == 0, don't do this check.
			if tokenInMaxs.Len() != 0 {
				if !(neededLpLiquidity.DenomsSubsetOf(tokenInMaxs)) {
					return nil, sdkmath.ZeroInt(), errorsmod.Wrapf(types.ErrLimitMaxAmount, "TokenInMaxs does not include all the tokens that are part of the target pool,"+
						" upperbound: %v, needed %v", tokenInMaxs, neededLpLiquidity)
				} else if !(tokenInMaxs.DenomsSubsetOf(neededLpLiquidity)) {
					return nil, sdkmath.ZeroInt(), errorsmod.Wrapf(types.ErrDenomNotFoundInPool, "TokenInMaxs includes tokens that are not part of the target pool,"+
						" input tokens: %v, pool tokens %v", tokenInMaxs, neededLpLiquidity)
				}
				if !(tokenInMaxs.IsAllGTE(neededLpLiquidity)) {
					return nil, sdkmath.ZeroInt(), errorsmod.Wrapf(types.ErrLimitMaxAmount, "TokenInMaxs is less than the needed LP liquidity to this JoinPoolNoSwap,"+
						" upperbound: %v, needed %v", tokenInMaxs, neededLpLiquidity)
				}
			}

			tokensIn = neededLpLiquidity
		}
		params := k.GetParams(ctx)
		snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
		tokensJoined, sharesOut, slippage, weightBalanceBonus, swapFee, err := pool.JoinPool(ctx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokensIn, params)
		if err != nil {
			return nil, sdkmath.ZeroInt(), err
		}

		if pool.PoolParams.UseOracle && len(tokensIn) == 1 {
			slippageAmount := slippage.Mul(tokensIn[0].Amount.ToLegacyDec()).RoundInt()
			k.TrackWeightBreakingSlippage(ctx, pool.PoolId, sdk.NewCoin(tokensIn[0].Denom, slippageAmount))
		}

		// sanity check, don't return error as not worth halting the LP. We know its not too much.
		if sharesOut.LT(shareOutAmount) {
			ctx.Logger().Error(fmt.Sprintf("Expected to JoinPoolNoSwap >= %s shares, actually did %s shares",
				shareOutAmount, sharesOut))
		}

		err = k.ApplyJoinPoolStateChange(ctx, pool, sender, sharesOut, tokensJoined, weightBalanceBonus, swapFee)
		if err != nil {
			return nil, sdkmath.Int{}, err
		}
		// Increase liquidity amount
		err = k.RecordTotalLiquidityIncrease(ctx, tokensJoined)
		if err != nil {
			return nil, sdkmath.Int{}, err
		}

		return tokensJoined, sharesOut, err
	}

	params := k.GetParams(ctx)
	// on oracle pool, full tokenInMaxs are used regardless shareOutAmount
	snapshot := k.GetAccountedPoolSnapshotOrSet(ctx, pool)
	tokensJoined, sharesOut, _, weightBalanceBonus, swapFee, err := pool.JoinPool(ctx, &snapshot, k.oracleKeeper, k.accountedPoolKeeper, tokenInMaxs, params)
	if err != nil {
		return nil, sdkmath.ZeroInt(), err
	}

	// sanity check, don't return error as not worth halting the LP. We know its not too much.
	if sharesOut.LT(shareOutAmount) {
		ctx.Logger().Error(fmt.Sprintf("Expected to JoinPoolNoSwap >= %s shares, actually did %s shares",
			shareOutAmount, sharesOut))
	}

	err = k.ApplyJoinPoolStateChange(ctx, pool, sender, sharesOut, tokensJoined, weightBalanceBonus, swapFee)
	if err != nil {
		return nil, sdkmath.Int{}, err
	}

	// Increase liquidity amount
	err = k.RecordTotalLiquidityIncrease(ctx, tokensJoined)
	if err != nil {
		return nil, sdkmath.Int{}, err
	}

	return tokensJoined, sharesOut, err
}
