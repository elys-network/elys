package keeper

import (
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
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
	swapFeeIn sdk.Dec,
	swapFeeOut sdk.Dec,
	weightBalanceBonus sdk.Dec,
) error {
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	err := k.SetPool(ctx, pool)
	if err != nil {
		return err
	}

	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err = k.bankKeeper.SendCoins(ctx, sender, poolAddr, tokensIn)
	if err != nil {
		return err
	}

	swapFeeInCoins := portionCoins(tokensIn, swapFeeIn)
	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
		if err != nil {
			return err
		}
		k.OnCollectFee(ctx, pool, swapFeeInCoins)
	}

	err = k.bankKeeper.SendCoins(ctx, poolAddr, sender, sdk.Coins{tokenOut})
	if err != nil {
		return err
	}

	swapFeeOutCoins := portionCoins(tokensOut, swapFeeOut)
	if swapFeeOutCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, sender, rebalanceTreasury, swapFeeOutCoins)
		if err != nil {
			return err
		}
		k.OnCollectFee(ctx, pool, swapFeeOutCoins)
	}

	// calculate treasury amount to send as bonus
	rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
	treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount
	bonusTokenAmount := sdk.NewDecFromInt(tokenOut.Amount).Mul(weightBalanceBonus).RoundInt()
	if treasuryTokenAmount.LT(bonusTokenAmount) {
		bonusTokenAmount = treasuryTokenAmount
	}

	// send bonus tokens to sender if positive
	if weightBalanceBonus.IsPositive() && bonusTokenAmount.IsPositive() {
		bonusToken := sdk.NewCoin(tokenOut.Denom, bonusTokenAmount)
		err = k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, sender, sdk.Coins{bonusToken})
		if err != nil {
			return err
		}
	}

	types.EmitSwapEvent(ctx, sender, pool.GetPoolId(), tokensIn, tokensOut)
	k.hooks.AfterSwap(ctx, sender, pool.GetPoolId(), tokensIn, tokensOut)
	k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	k.RecordTotalLiquidityDecrease(ctx, tokensOut)
	k.RecordTotalLiquidityDecrease(ctx, swapFeeInCoins)

	return nil
}
