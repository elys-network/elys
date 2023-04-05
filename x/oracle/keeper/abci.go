package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) {
	// Remove outdated prices
	params := k.GetParams(ctx)
	for _, price := range k.GetAllPrice(ctx) {
		if price.Timestamp+params.PriceExpiryTime < uint64(ctx.BlockTime().Unix()) {
			k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
		}
	}
}
