package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) SetPoolInfo(ctx sdk.Context, pool types.PoolInfo) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolInfoKey(pool.GetPoolId())
	b := k.cdc.MustMarshal(&pool)
	store.Set(key, b)
}

func (k Keeper) GetPoolInfo(ctx sdk.Context, poolId uint64) (val types.PoolInfo, found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolInfoKey(poolId)

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemovePoolInfo(ctx sdk.Context, poolId uint64) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPoolInfoKey(poolId)
	store.Delete(key)
}

func (k Keeper) GetAllPoolInfos(ctx sdk.Context) (list []types.PoolInfo) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PoolInfoKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PoolInfo
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) RemoveLegacyPoolInfo(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyPoolInfoKeyPrefix))
	store.Delete(types.LegacyPoolInfoKey(poolId))
}

func (k Keeper) GetAllLegacyPoolInfos(ctx sdk.Context) (list []types.PoolInfo) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyPoolInfoKeyPrefix))
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
		p, found := k.GetPoolInfo(ctx, pm.PoolId)
		if found {
			p.Multiplier = pm.Multiplier
			k.SetPoolInfo(ctx, p)

			ctx.EventManager().EmitEvents(sdk.Events{
				sdk.NewEvent(
					types.TypeEvtSetPoolMultiplier,
					sdk.NewAttribute(types.AttributePoolId, fmt.Sprintf("%d", pm.PoolId)),
					sdk.NewAttribute(types.AttributeMultiplier, pm.Multiplier.String()),
				),
			})
		}
	}

	return true
}
