package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
)

func (k Keeper) GetExternalIncentiveIndex(ctx sdk.Context) (index uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveIndexKeyPrefix))
	index = sdk.BigEndianToUint64(store.Get(types.ExternalIncentiveIndex()))
	return index
}

func (k Keeper) SetExternalIncentiveIndex(ctx sdk.Context, index uint64) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveIndexKeyPrefix))
	store.Set(types.ExternalIncentiveIndex(), sdk.Uint64ToBigEndian(index))
}

func (k Keeper) SetExternalIncentive(ctx sdk.Context, externalIncentive types.ExternalIncentive) error {
	// Update external incentive index and increase +1
	index := k.GetExternalIncentiveIndex(ctx)
	externalIncentive.Id = index
	k.SetExternalIncentiveIndex(ctx, index+1)

	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ExternalIncentiveKeyPrefix))
	b := k.cdc.MustMarshal(&externalIncentive)
	store.Set(types.ExternalIncentiveKey(externalIncentive.Id), b)
	return nil
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
