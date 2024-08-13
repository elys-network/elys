package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) GetExternalIncentiveIndex(ctx sdk.Context) (index uint64) {
	store := ctx.KVStore(k.storeKey)
	v := store.Get(types.ExternalIncentiveIndexKeyPrefix)
	index = sdk.BigEndianToUint64(v)
	return index
}

func (k Keeper) SetExternalIncentiveIndex(ctx sdk.Context, index uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.ExternalIncentiveIndexKeyPrefix, sdk.Uint64ToBigEndian(index))
}

// remove after migration
func (k Keeper) GetLegacyExternalIncentiveIndex(ctx sdk.Context) (index uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyExternalIncentiveIndexKeyPrefix))
	index = sdk.BigEndianToUint64(store.Get(types.LegacyExternalIncentiveIndex()))
	return index
}

// remove after migration
func (k Keeper) RemoveLegacyExternalIncentiveIndex(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyExternalIncentiveIndexKeyPrefix))
	store.Delete(types.LegacyExternalIncentiveIndex())
}

func (k Keeper) SetExternalIncentive(ctx sdk.Context, externalIncentive types.ExternalIncentive) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveKeyPrefix))
	b := k.cdc.MustMarshal(&externalIncentive)
	store.Set(types.ExternalIncentiveKey(externalIncentive.Id), b)
}

func (k Keeper) GetExternalIncentive(ctx sdk.Context, id uint64) (val types.ExternalIncentive, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveKeyPrefix))
	b := store.Get(types.ExternalIncentiveKey(id))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) RemoveExternalIncentive(ctx sdk.Context, id uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveKeyPrefix))
	store.Delete(types.ExternalIncentiveKey(id))
}

func (k Keeper) GetAllExternalIncentives(ctx sdk.Context) (list []types.ExternalIncentive) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.ExternalIncentive
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
