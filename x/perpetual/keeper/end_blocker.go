package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) EndBlocker(ctx sdk.Context) {
	pools := k.GetAllPools(ctx)

	for _, pool := range pools {
		cacheCtx, write := ctx.CacheContext()
		err := k.ClosePositionsOnADL(cacheCtx, pool)
		if err == nil {
			write()
		}
	}
}
