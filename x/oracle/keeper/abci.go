package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) {
	// Remove outdated prices
	params := k.GetParams(ctx)
	for _, price := range k.GetAllLegacyPrice(ctx) {
		if price.Timestamp+params.PriceExpiryTime < uint64(ctx.BlockTime().Unix()) || price.BlockHeight+params.LifeTimeInBlocks < uint64(ctx.BlockHeight()) {
			k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
		}
	}
}
