package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/masterchef/types"
)

func (k Keeper) GetExternalIncentiveIndex(ctx sdk.Context) (index uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	v := store.Get(types.ExternalIncentiveIndexKeyPrefix)
	index = sdk.BigEndianToUint64(v)
	return index
}

func (k Keeper) SetExternalIncentiveIndex(ctx sdk.Context, index uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.ExternalIncentiveIndexKeyPrefix, sdk.Uint64ToBigEndian(index))
}

func (k Keeper) SetExternalIncentive(ctx sdk.Context, externalIncentive types.ExternalIncentive) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetExternalIncentiveKey(externalIncentive.Id)
	b := k.cdc.MustMarshal(&externalIncentive)
	store.Set(key, b)
}

func (k Keeper) GetExternalIncentive(ctx sdk.Context, id uint64) (val types.ExternalIncentive, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetExternalIncentiveKey(id)
	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveExternalIncentive(ctx sdk.Context, id uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetExternalIncentiveKey(id)
	store.Delete(key)
}

func (k Keeper) GetAllExternalIncentives(ctx sdk.Context) (list []types.ExternalIncentive) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.ExternalIncentiveKeyPrefix)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ExternalIncentive
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
