package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const OracleLifeTimeInBlocks = uint64(1)

func (k Keeper) EndBlock(ctx sdk.Context) {
	// Remove outdated prices
	params := k.GetParams(ctx)
	for _, price := range k.GetAllPrice(ctx) {
		if price.Timestamp+params.PriceExpiryTime < uint64(ctx.BlockTime().Unix()) {
			k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
		}

		if price.BlockHeight+OracleLifeTimeInBlocks < uint64(ctx.BlockHeight()) {
			k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
		}
	}
}
