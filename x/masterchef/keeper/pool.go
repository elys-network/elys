package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetPool(ctx sdk.Context, pool types.PoolInfo) error {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolInfoKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolInfoKey(pool.PoolId), b)
	return nil
}

func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.PoolInfo, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolInfoKeyPrefix))

	b := store.Get(types.PoolInfoKey(poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePool(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolInfoKeyPrefix))
	store.Delete(types.PoolInfoKey(poolId))
}

func (k Keeper) GetAllPools(ctx sdk.Context) (list []types.PoolInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PoolInfoKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) UpdatePoolMultipliers(ctx sdk.Context, poolMultipliers []types.PoolMultiplier) bool {
	if len(poolMultipliers) < 1 {
		return false
	}

	// Update pool multiplier
	for _, pm := range poolMultipliers {
		p, found := k.GetPool(ctx, pm.PoolId)
		if found {
			p.Multiplier = pm.Multiplier
			k.SetPool(ctx, p)
		}
	}

	return true
}
