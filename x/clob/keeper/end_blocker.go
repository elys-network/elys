package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) error {
	//allMarkets := k.GetAllPerpetualMarket(ctx)
	//
	//for _, market := range allMarkets {
	//	cacheCtx, writeCache := ctx.CacheContext()
	//
	//	_, _, _, err := k.ExecuteMarket(cacheCtx, market.Id, 20)
	//	if err != nil {
	//		ctx.Logger().Error("Failed to execute market", "market_id", market.Id, "error", err)
	//	} else {
	//		writeCache()
	//	}
	//}
	return nil
}
