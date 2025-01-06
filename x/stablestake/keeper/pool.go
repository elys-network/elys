package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

// GetPools get pool as types.Pool
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (pool types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetPoolKey(poolId))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &pool)
	return
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.GetPoolKey(pool.PoolId), b)
}

func (k Keeper) GetRedemptionRateForPool(ctx sdk.Context, pool types.Pool) math.LegacyDec {
	totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForPool(pool.PoolId))

	if totalShares.Amount.IsZero() {
		return math.LegacyZeroDec()
	}

	return pool.TotalValue.ToLegacyDec().Quo(totalShares.Amount.ToLegacyDec())
}
