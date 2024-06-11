<!--
order: 3
-->

# Keeper

## Price Management

The `oracle` module's keeper handles the management and querying of asset prices and related information. It ensures the timely update and retrieval of price data and manages the lifecycle of price feeders.

### EndBlocker

The `EndBlocker` function is invoked at the end of each block. It is responsible for removing outdated prices based on the configured expiration parameters.

```go
func (k Keeper) EndBlock(ctx sdk.Context) {
    params := k.GetParams(ctx)
    for _, price := range k.GetAllPrice(ctx) {
        if price.Timestamp + params.PriceExpiryTime < uint64(ctx.BlockTime().Unix()) {
            k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
        }
        if price.BlockHeight + params.LifeTimeInBlocks < uint64(ctx.BlockHeight()) {
            k.RemovePrice(ctx, price.Asset, price.Source, price.Timestamp)
        }
    }
}
```

### Managing Asset Info

The `SetAssetInfo`, `GetAssetInfo`, `RemoveAssetInfo`, and `GetAllAssetInfo` functions handle the creation, retrieval, deletion, and listing of asset information.

```go
func (k Keeper) SetAssetInfo(ctx sdk.Context, assetInfo types.AssetInfo) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
    bz := k.cdc.MustMarshal(&assetInfo)
    store.Set(types.AssetInfoKey(assetInfo.Denom), bz)
}

func (k Keeper) GetAssetInfo(ctx sdk.Context, denom string) (val types.AssetInfo, found bool) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
    bz := store.Get(types.AssetInfoKey(denom))
    if bz == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(bz, &val)
    return val, true
}

func (k Keeper) RemoveAssetInfo(ctx sdk.Context, denom string) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
    store.Delete(types.AssetInfoKey(denom))
}

func (k Keeper) GetAllAssetInfo(ctx sdk.Context) (list []types.AssetInfo) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.AssetInfoKeyPrefix))
    iterator := sdk.KVStorePrefixIterator(store, []byte{})
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var val types.AssetInfo
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }
    return
}
```

### Managing Band Price Results

The `SetBandPriceResult`, `GetBandPriceResult`, `GetLastBandRequestId`, and `SetLastBandRequestId` functions manage the storage and retrieval of price data from Band protocol.

```go
func (k Keeper) SetBandPriceResult(ctx sdk.Context, requestID types.OracleRequestID, result types.BandPriceResult) {
    store := ctx.KVStore(k.storeKey)
    store.Set(types.BandPriceResultStoreKey(requestID), k.cdc.MustMarshal(&result))
}

func (k Keeper) GetBandPriceResult(ctx sdk.Context, id types.OracleRequestID) (types.BandPriceResult, error) {
    bz := ctx.KVStore(k.storeKey).Get(types.BandPriceResultStoreKey(id))
    if bz == nil {
        return types.BandPriceResult{}, errorsmod.Wrapf(types.ErrNotAvailable, "Result for request ID %d is not available.", id)
    }
    var result types.BandPriceResult
    k.cdc.MustUnmarshal(bz, &result)
    return result, nil
}

func (k Keeper) GetLastBandRequestId(ctx sdk.Context) int64 {
    bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastBandRequestIdKey))
    intV := gogotypes.Int64Value{}
    k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
    return intV.GetValue()
}

func (k Keeper) SetLastBandRequestId(ctx sdk.Context, id types.OracleRequestID) {
    store := ctx.KVStore(k.storeKey)
    store.Set(types.KeyPrefix(types.LastBandRequestIdKey), k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}
```

### Managing Prices

The `SetPrice`, `GetPrice`, `GetLatestPriceFromAssetAndSource`, `GetLatestPriceFromAnySource`, `RemovePrice`, and `GetAllPrice` functions manage the lifecycle and retrieval of price data.

```go
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
    store := ctx.KVStore(k.storeKey)
    b := k.cdc.MustMarshal(&price)
    store.Set(types.PriceKey(price.Asset, price.Source, price.Timestamp), b)
}

func (k Keeper) GetPrice(ctx sdk.Context, asset, source string, timestamp uint64) (val types.Price, found bool) {
    store := ctx.KVStore(k.storeKey)
    b := store.Get(types.PriceKey(asset, source, timestamp))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}

func (k Keeper) GetLatestPriceFromAssetAndSource(ctx sdk.Context, asset, source string) (val types.Price, found bool) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStoreReversePrefixIterator(store, types.PriceKeyPrefixAssetAndSource(asset, source))
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var val types.Price
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        return val, true
    }
    return val, false
}

func (k Keeper) GetLatestPriceFromAnySource(ctx sdk.Context, asset string) (val types.Price, found bool) {
    store := ctx.KVStore(k.storeKey)
    iterator := sdk.KVStoreReversePrefixIterator(store, types.PriceKeyPrefixAsset(asset))
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var val types.Price
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        return val, true
    }
    return val, false
}

func (k Keeper) RemovePrice(ctx sdk.Context, asset, source string, timestamp uint64) {
    store := ctx.KVStore(k.storeKey)
    store.Delete(types.PriceKey(asset, source, timestamp))
}

func (k Keeper) GetAllPrice(ctx sdk.Context) (list []types.Price) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceKeyPrefix))
    iterator := sdk.KVStorePrefixIterator(store, []byte{})
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var val types.Price
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }
    return
}
```

### Managing Price Feeders

The `SetPriceFeeder`, `GetPriceFeeder`, `RemovePriceFeeder`, and `GetAllPriceFeeder` functions manage the lifecycle and retrieval of price feeder data.

```go
func (k Keeper) SetPriceFeeder(ctx sdk.Context, priceFeeder types.PriceFeeder) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceFeederKeyPrefix))
    b := k.cdc.MustMarshal(&priceFeeder)
    store.Set(types.PriceFeederKey(priceFeeder.Feeder), b)
}

func (

k Keeper) GetPriceFeeder(ctx sdk.Context, feeder string) (val types.PriceFeeder, found bool) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceFeederKeyPrefix))
    b := store.Get(types.PriceFeederKey(feeder))
    if b == nil {
        return val, false
    }
    k.cdc.MustUnmarshal(b, &val)
    return val, true
}

func (k Keeper) RemovePriceFeeder(ctx sdk.Context, feeder string) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceFeederKeyPrefix))
    store.Delete(types.PriceFeederKey(feeder))
}

func (k Keeper) GetAllPriceFeeder(ctx sdk.Context) (list []types.PriceFeeder) {
    store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceFeederKeyPrefix))
    iterator := sdk.KVStorePrefixIterator(store, []byte{})
    defer iterator.Close()
    for ; iterator.Valid(); iterator.Next() {
        var val types.PriceFeeder
        k.cdc.MustUnmarshal(iterator.Value(), &val)
        list = append(list, val)
    }
    return
}
```
