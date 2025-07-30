package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlock(ctx sdk.Context) {
	for _, info := range k.GetAllAssetInfo(ctx) {
		prices := k.GetAllAssetPrice(ctx, info.Display)
		if len(prices) <= 1 {
			continue // nothing to prune
		}

		// Find the newest price
		latestIdx, latestTs := 0, prices[0].Timestamp
		for i := 1; i < len(prices); i++ {
			if prices[i].Timestamp > latestTs {
				latestIdx, latestTs = i, prices[i].Timestamp
			}
		}

		// Remove everything except the newest
		for i, p := range prices {
			if i == latestIdx {
				continue
			}
			k.RemovePrice(ctx, p.Asset, p.Source, p.Timestamp)
		}
	}
}
