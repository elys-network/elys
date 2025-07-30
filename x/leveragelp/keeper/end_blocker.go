package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) EndBlocker(ctx sdk.Context) {
	allPools := k.GetAllPools(ctx)
	for _, pool := range allPools {
		cacheCtx, write := ctx.CacheContext()
		err := k.ClosePositionsOnADL(cacheCtx, pool)
		if err == nil {
			write()
		}
	}
}
