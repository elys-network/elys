package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

// SetEntry set a specific entry in the store from its index
func (k Keeper) SetEntry(ctx sdk.Context, entry types.Entry) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
	b := k.cdc.MustMarshal(&entry)
	store.Set(types.EntryKey(entry.BaseDenom), b)
}

// GetEntry returns a entry from its index
func (k Keeper) GetEntry(ctx sdk.Context, baseDenom string) (val types.Entry, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
	b := store.Get(types.EntryKey(baseDenom))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// GetEntryByDenom returns a entry from its denom value
func (k Keeper) GetEntryByDenom(ctx sdk.Context, denom string) (val types.Entry, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Entry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.Denom == denom {
			return val, true
		}
	}

	return types.Entry{}, false
}

// RemoveEntry removes a entry from the store
func (k Keeper) RemoveEntry(ctx sdk.Context, baseDenom string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
	store.Delete(types.EntryKey(baseDenom))
}

// GetAllEntry returns all entry
func (k Keeper) GetAllEntry(ctx sdk.Context) (list []types.Entry) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Entry
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetUsdcDenom(ctx sdk.Context) (string, bool) {
	entry, found := k.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return ptypes.BaseCurrency, false
	}
	return entry.Denom, true
}
