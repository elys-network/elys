<!--
order: 3
-->

# Keeper

## Asset Management

The `assetprofile` module's keeper handles asset profile management, including creating, updating, and deleting asset entries. It ensures that asset properties are properly defined and managed within the network.

### Setting an Asset Entry

The `SetEntry` function stores a new asset entry or updates an existing one.

```go
func (k Keeper) SetEntry(ctx sdk.Context, entry types.Entry) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    b := k.cdc.MustMarshal(&entry)
    store.Set(types.EntryKey(entry.BaseDenom), b)
}
```

### Getting an Asset Entry

The `GetEntry` function retrieves an asset entry based on its base denomination.

```go
func (k Keeper) GetEntry(ctx sdk.Context, baseDenom string) (val types.Entry, found bool) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    b := store.Get(types.EntryKey(baseDenom))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}
```

### Removing an Asset Entry

The `RemoveEntry` function deletes an asset entry from the store.

```go
func (k Keeper) RemoveEntry(ctx sdk.Context, baseDenom string) {
    store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.EntryKeyPrefix))
    store.Delete(types.EntryKey(baseDenom))
}
```

### Retrieving All Asset Entries

The `GetAllEntry` function retrieves all asset entries from the store.

```go
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
```
