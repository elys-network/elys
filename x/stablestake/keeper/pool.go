package keeper

import (
	"cosmossdk.io/math"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetAllPools(ctx sdk.Context) (pools []types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, types.DebtPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pool := types.Pool{}
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		pools = append(pools, pool)
	}

	return
}

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

func (k Keeper) GetPoolByDenom(ctx sdk.Context, denom string) (types.Pool, bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	iterator := storetypes.KVStorePrefixIterator(store, types.DebtPrefixKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		pool := types.Pool{}
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if pool.DepositDenom == denom {
			return pool, true
		}
	}
	return types.Pool{}, false
}

func (k Keeper) GetRedemptionRateForPool(ctx sdk.Context, pool types.Pool) math.LegacyDec {
	totalShares := k.bk.GetSupply(ctx, types.GetShareDenomForPool(pool.PoolId))

	if totalShares.Amount.IsZero() {
		return math.LegacyZeroDec()
	}

	return pool.TotalValue.ToLegacyDec().Quo(totalShares.Amount.ToLegacyDec())
}
