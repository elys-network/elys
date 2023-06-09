package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// updatePoolForSwap takes a pool, sender, and tokenIn, tokenOut amounts
// It then updates the pool's balances to the new reserve amounts, and
// sends the in tokens from the sender to the pool, and the out tokens from the pool to the sender.
func (k Keeper) updatePoolForSwap(
	ctx sdk.Context,
	pool types.Pool,
	sender sdk.AccAddress,
	tokenIn sdk.Coin,
	tokenOut sdk.Coin,
) error {
	tokensIn := sdk.Coins{tokenIn}
	tokensOut := sdk.Coins{tokenOut}

	err := k.SetPool(ctx, pool)
	if err != nil {
		return err
	}

	poolAddr := sdk.MustAccAddressFromBech32(pool.GetAddress())
	err = k.bankKeeper.SendCoins(ctx, sender, poolAddr, sdk.Coins{tokenIn})
	if err != nil {
		return err
	}

	swapFeeCoins := portionCoins(sdk.Coins{tokenIn}, pool.PoolParams.SwapFee)
	rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
	err = k.bankKeeper.SendCoins(ctx, poolAddr, rebalanceTreasury, swapFeeCoins)
	if err != nil {
		return err
	}
	k.OnCollectFee(ctx, pool, swapFeeCoins)

	err = k.bankKeeper.SendCoins(ctx, poolAddr, sender, sdk.Coins{tokenOut})
	if err != nil {
		return err
	}

	types.EmitSwapEvent(ctx, sender, pool.GetPoolId(), tokensIn, tokensOut)
	k.hooks.AfterSwap(ctx, sender, pool.GetPoolId(), tokensIn, tokensOut)
	k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	k.RecordTotalLiquidityDecrease(ctx, tokensOut)

	return err
}
