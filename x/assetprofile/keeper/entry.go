package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/assetprofile/types"
)

// SetEntry set a specific entry in the store from its index
func (k Keeper) SetEntry(ctx sdk.Context, entry types.Entry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EntryKeyPrefix))
	b := k.cdc.MustMarshal(&entry)
	store.Set(types.EntryKey(entry.BaseDenom), b)
}

// GetEntry returns a entry from its index
func (k Keeper) GetEntry(ctx sdk.Context, baseDenom string) (val types.Entry, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EntryKeyPrefix))
	b := store.Get(types.EntryKey(baseDenom))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveEntry removes a entry from the store
func (k Keeper) RemoveEntry(ctx sdk.Context, baseDenom string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EntryKeyPrefix))
	store.Delete(types.EntryKey(baseDenom))
}

// GetAllEntry returns all entry
func (k Keeper) GetAllEntry(ctx sdk.Context) (list []types.Entry) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.EntryKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Entry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
