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

	acc := sdk.AccAddress(pool.GetAddress())

	err = k.bankKeeper.SendCoins(ctx, sender, acc, sdk.Coins{
		tokenIn,
	})
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoins(ctx, acc, sender, sdk.Coins{
		tokenOut,
	})
	if err != nil {
		return err
	}

	types.EmitSwapEvent(ctx, sender, pool.GetPoolId(), tokensIn, tokensOut)
	k.hooks.AfterSwap(ctx, sender, pool.GetPoolId(), tokensIn, tokensOut)
	k.RecordTotalLiquidityIncrease(ctx, tokensIn)
	k.RecordTotalLiquidityDecrease(ctx, tokensOut)

	return err
}
