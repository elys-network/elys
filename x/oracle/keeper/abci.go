package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"sort"
)

func (k Keeper) EndBlock(ctx sdk.Context) {

	assetInfos := k.GetAllAssetInfo(ctx)

	for _, info := range assetInfos {
		allAssetPrice := k.GetAllAssetPrice(ctx, info.Display)
		total := len(allAssetPrice)

		// Need to sort it because order fetched from will not be ascending order depending on source
		// If we remove the source then this should not be needed
		sort.Slice(allAssetPrice, func(i, j int) bool {
			return allAssetPrice[i].Timestamp < allAssetPrice[j].Timestamp
		})

		for i, price := range allAssetPrice {
			// We don't remove the last element
			if i < total-1 {
				k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
			}
		}
	}
}
