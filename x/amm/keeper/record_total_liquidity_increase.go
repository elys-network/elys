package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) RecordTotalLiquidityIncrease(ctx sdk.Context, coins sdk.Coins) {
	for _, coin := range coins {
		amount := k.GetDenomLiquidity(ctx, coin.Denom)
		amount = amount.Add(coin.Amount)
		k.setDenomLiquidity(ctx, coin.Denom, amount)
	}
}
