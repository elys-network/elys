package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// UpdatePoolForSwap takes a pool, sender, and tokenIn, tokenOut amounts
// It then updates the pool's balances to the new reserve amounts, and
// sends the in tokens from the sender to the pool, and the out tokens from the pool to the sender.

func (k Keeper) UpdatePoolForSwap(
	ctx sdk.Context,
	pool types.Pool,
	sender sdk.AccAddress,
	recipient sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
	swapFeeIn sdkmath.LegacyDec,
	swapFeeOut sdkmath.LegacyDec,
	weightBalanceBonus sdkmath.LegacyDec,
) (sdkmath.Int, error) {
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	k.SetPool(ctx, pool)

	// send tokensIn from sender to pool
	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err := k.bankKeeper.SendCoins(ctx, sender, poolAddr, tokensIn)
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeInCoins := sdk.Coins{}
	if !weightBalanceBonus.IsPositive() {
		swapFeeInCoins = PortionCoins(tokensIn, swapFeeIn)
	}

	// send swap fee to rebalance treasury
	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
		err = k.OnCollectFee(ctx, pool, swapFeeInCoins)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
	}

	// init weightRecoveryFeeAmount to zero
	weightRecoveryFeeAmount := sdkmath.ZeroInt()

	// send half (weight breaking fee portion) of weight breaking fee to rebalance treasury
	if !weightBalanceBonus.IsPositive() {
		params := k.GetParams(ctx)
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		weightRecoveryFee := weightBalanceBonus.Abs().Mul(params.WeightBreakingFeePortion)

		// weight recovery fee must be positive
		if weightRecoveryFee.IsPositive() {
			weightRecoveryFeeAmount = tokenOut.Amount.ToLegacyDec().Mul(weightRecoveryFee).RoundInt()

			// send weight recovery fee to rebalance treasury if weight recovery fee amount is positive
			if weightRecoveryFeeAmount.IsPositive() {
				err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, sdk.Coins{sdk.NewCoin(tokenOut.Denom, weightRecoveryFeeAmount)})
				if err != nil {
					return sdkmath.ZeroInt(), err
				}
			}
		}
	}

	// Send coins to recipient
	err = k.bankKeeper.SendCoins(ctx, poolAddr, recipient, sdk.Coins{tokenOut})
	if err != nil {
		return sdkmath.ZeroInt(), err
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeOutCoins := sdk.Coins{}
	if !weightBalanceBonus.IsPositive() {
		swapFeeOutCoins = PortionCoins(tokensOut, swapFeeOut)
	}
	if swapFeeOutCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeOutCoins)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
		err = k.OnCollectFee(ctx, pool, swapFeeOutCoins)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
	}

	// calculate total swap fee
	swapFeeCoins := swapFeeInCoins.Add(swapFeeOutCoins...)

	// calculate bonus token amount if weightBalanceBonus is positive
	if weightBalanceBonus.IsPositive() {
		// get treasury balance
		rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount

		// bonus token amount is the tokenOut amount times weightBalanceBonus
		bonusTokenAmount := sdkmath.LegacyNewDecFromInt(tokenOut.Amount).Mul(weightBalanceBonus).RoundInt()

		// if treasury balance is less than bonusTokenAmount, set bonusTokenAmount to treasury balance
		if treasuryTokenAmount.LT(bonusTokenAmount) {
			bonusTokenAmount = treasuryTokenAmount
		}

		// send bonusTokenAmount from pool addr to recipient addr, we are shortcutting the rebalance treasury address to optimize gas
		if bonusTokenAmount.IsPositive() {
			bonusToken := sdk.NewCoin(tokenOut.Denom, bonusTokenAmount)
			err = k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, recipient, sdk.Coins{bonusToken})
			if err != nil {
				return sdkmath.ZeroInt(), err
			}
		}
	}

	// emit swap event
	types.EmitSwapEvent(ctx, sender, recipient, pool.GetPoolId(), tokensIn, tokensOut)
	if k.hooks != nil {
		err = k.hooks.AfterSwap(ctx, sender, pool, tokensIn, tokensOut)
		if err != nil {
			return sdkmath.ZeroInt(), err
		}
	}

	// record tokenIn amount as total liquidity increase
	err = k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	if err != nil {
		return sdkmath.Int{}, err
	}

	// record tokenOut amount as total liquidity decrease
	err = k.RecordTotalLiquidityDecrease(ctx, tokensOut)
	if err != nil {
		return sdkmath.Int{}, err
	}

	// record swap fee as total liquidity decrease
	err = k.RecordTotalLiquidityDecrease(ctx, swapFeeCoins)
	if err != nil {
		return sdkmath.Int{}, err
	}

	// record weight recovery fee amount as total liquidity decrease
	weightRecoveryToken := sdk.NewCoin(tokenOut.Denom, weightRecoveryFeeAmount)
	err = k.RecordTotalLiquidityDecrease(ctx, sdk.Coins{weightRecoveryToken})
	if err != nil {
		return sdkmath.Int{}, err
	}

	// return swap fee out amount
	return swapFeeOutCoins.AmountOf(tokenOut.Denom), nil
}
