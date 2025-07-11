package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) EndBlocker(ctx sdk.Context) {
	allPools := k.GetAllPools(ctx)
	for _, pool := range allPools {
		ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
		if err != nil {
			return
		}
		currentLeverageRatio := pool.LeveragedLpAmount.ToLegacyDec().Quo(ammPool.TotalShares.Amount.ToLegacyDec())

		if currentLeverageRatio.GT(pool.AdlTriggerRatio) {
			cacheCtx, write := ctx.CacheContext()
			err = k.ClosePositionsOnADL(cacheCtx, pool, currentLeverageRatio)
			if err == nil {
				write()
			}
		}
	}
}
