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
	recipient sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
	swapFeeIn sdk.Dec,
	swapFeeOut sdk.Dec,
	weightBalanceBonus sdk.Dec,
) (error, sdk.Int) {
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	err := k.SetPool(ctx, pool)
	if err != nil {
		return err, sdk.ZeroInt()
	}

	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err = k.bankKeeper.SendCoins(ctx, sender, poolAddr, tokensIn)
	if err != nil {
		return err, sdk.ZeroInt()
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeInCoins := sdk.Coins{}
	if weightBalanceBonus.IsZero() {
		swapFeeInCoins = PortionCoins(tokensIn, swapFeeIn)
	}

	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
		if err != nil {
			return err, sdk.ZeroInt()
		}
		err = k.OnCollectFee(ctx, pool, swapFeeInCoins)
		if err != nil {
			return err, sdk.ZeroInt()
		}
	}

	// Send coins to recipient
	err = k.bankKeeper.SendCoins(ctx, poolAddr, recipient, sdk.Coins{tokenOut})
	if err != nil {
		return err, sdk.ZeroInt()
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeOutCoins := sdk.Coins{}
	if weightBalanceBonus.IsZero() {
		swapFeeOutCoins = PortionCoins(tokensOut, swapFeeOut)
	}
	if swapFeeOutCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, recipient, rebalanceTreasury, swapFeeOutCoins)
		if err != nil {
			return err, sdk.ZeroInt()
		}
		err = k.OnCollectFee(ctx, pool, swapFeeOutCoins)
		if err != nil {
			return err, sdk.ZeroInt()
		}
	}

	swapFeeCoins := swapFeeInCoins.Add(swapFeeOutCoins...)

	// calculate treasury amount to send as bonus
	rebalanceTreasuryAddr := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
	treasuryTokenAmount := k.bankKeeper.GetBalance(ctx, rebalanceTreasuryAddr, tokenOut.Denom).Amount
	bonusTokenAmount := sdk.NewDecFromInt(tokenOut.Amount).Mul(weightBalanceBonus).RoundInt()
	if treasuryTokenAmount.LT(bonusTokenAmount) {
		bonusTokenAmount = treasuryTokenAmount
	}

	// send bonus tokens to recipient if positive
	if weightBalanceBonus.IsPositive() && bonusTokenAmount.IsPositive() {
		bonusToken := sdk.NewCoin(tokenOut.Denom, bonusTokenAmount)
		err = k.bankKeeper.SendCoins(ctx, rebalanceTreasuryAddr, recipient, sdk.Coins{bonusToken})
		if err != nil {
			return err, sdk.ZeroInt()
		}
	}

	types.EmitSwapEvent(ctx, sender, recipient, pool.GetPoolId(), tokensIn, tokensOut)
	if k.hooks != nil {
		k.hooks.AfterSwap(ctx, sender, pool, tokensIn, tokensOut)
	}
	k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	k.RecordTotalLiquidityDecrease(ctx, tokensOut)
	k.RecordTotalLiquidityDecrease(ctx, swapFeeCoins)

	return nil, swapFeeOutCoins.AmountOf(tokenOut.Denom)
}
