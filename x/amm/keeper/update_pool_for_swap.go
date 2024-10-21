package keeper

import (
	"cosmossdk.io/math"
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
) (math.Int, error) {
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	k.SetPool(ctx, pool)

	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err := k.bankKeeper.SendCoins(ctx, sender, poolAddr, tokensIn)
	if err != nil {
		return sdk.ZeroInt(), err
	}

	// apply swap fee when weight balance bonus is not available
	swapFeeInCoins := sdk.Coins{}
	if !weightBalanceBonus.IsPositive() {
		swapFeeInCoins = PortionCoins(tokensIn, swapFeeIn)
	}

	if swapFeeInCoins.IsAllPositive() {
		rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
		err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeInCoins)
		if err != nil {
			return sdk.ZeroInt(), err
		}
		err = k.OnCollectFee(ctx, pool, swapFeeInCoins)
		if err != nil {
			return sdk.ZeroInt(), err
		}
	}

	// Send coins to recipient
	err = k.bankKeeper.SendCoins(ctx, poolAddr, recipient, sdk.Coins{tokenOut})
	if err != nil {
		return sdk.ZeroInt(), err
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
			return sdk.ZeroInt(), err
		}
		err = k.OnCollectFee(ctx, pool, swapFeeOutCoins)
		if err != nil {
			return sdk.ZeroInt(), err
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
			return sdk.ZeroInt(), err
		}
	}

	types.EmitSwapEvent(ctx, sender, recipient, pool.GetPoolId(), tokensIn, tokensOut)
	if k.hooks != nil {
		err = k.hooks.AfterSwap(ctx, sender, pool, tokensIn, tokensOut)
		if err != nil {
			return math.ZeroInt(), err
		}
	}
	err = k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	if err != nil {
		return math.Int{}, err
	}
	err = k.RecordTotalLiquidityDecrease(ctx, tokensOut)
	if err != nil {
		return math.Int{}, err
	}
	err = k.RecordTotalLiquidityDecrease(ctx, swapFeeCoins)
	if err != nil {
		return math.Int{}, err
	}

	return swapFeeOutCoins.AmountOf(tokenOut.Denom), nil
}
