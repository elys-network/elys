package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// SetPrice set a specific price in the store from its index
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&price)
	store.Set(types.PriceKey(price.Asset, price.Source, price.Timestamp), b)
}

// GetPrice returns a price from its index
func (k Keeper) GetPrice(ctx sdk.Context, asset, source string, timestamp uint64) (val types.Price, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.PriceKey(asset, source, timestamp))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) GetLatestPriceFromAssetAndSource(ctx sdk.Context, asset, source string) (val types.Price, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStoreReversePrefixIterator(store, types.PriceKeyPrefixAssetAndSource(asset, source))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Price
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val, true
	}

	return val, false
}

func (k Keeper) GetLatestPriceFromAnySource(ctx sdk.Context, asset string) (val types.Price, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStoreReversePrefixIterator(store, types.PriceKeyPrefixAsset(asset))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Price
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		return val, true
	}

	return val, false
}

// RemovePrice removes a price from the store
func (k Keeper) RemovePrice(ctx sdk.Context, asset, source string, timestamp uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.PriceKey(asset, source, timestamp))
}

// GetAllPrice returns all price
func (k Keeper) GetAllPrice(ctx sdk.Context) (list []types.Price) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PriceKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Price
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// MigrateAllLegacyPrices migrates all legacy prices
func (k Keeper) MigrateAllLegacyPrices(ctx sdk.Context) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PriceKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyPrice
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		k.SetPrice(ctx, types.Price{
			Asset:       val.Asset,
			Price:       val.Price,
			Source:      val.Source,
			Provider:    val.Provider,
			Timestamp:   val.Timestamp,
			BlockHeight: uint64(ctx.BlockHeight()),
		})
	}

	return
}

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (types.Price, bool) {
	// try out elys source
	val, found := k.GetLatestPriceFromAssetAndSource(ctx, asset, types.ELYS)
	if found {
		return val, true
	}

	// try out band source
	val, found = k.GetLatestPriceFromAssetAndSource(ctx, asset, types.BAND)
	if found {
		return val, true
	}

	// find from any source if band source does not exist
	return k.GetLatestPriceFromAnySource(ctx, asset)
}

func Pow10(decimal uint64) (value sdkmath.LegacyDec) {
	value = sdkmath.LegacyNewDec(1)
	for i := 0; i < int(decimal); i++ {
		value = value.Mul(sdkmath.LegacyNewDec(10))
	}
	return
}

func (k Keeper) GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdkmath.LegacyDec {
	info, found := k.GetAssetInfo(ctx, denom)
	if !found {
		return sdkmath.LegacyZeroDec()
	}
	price, found := k.GetAssetPrice(ctx, info.Display)
	if !found {
		return sdkmath.LegacyZeroDec()
	}
	return price.Price.Quo(Pow10(info.Decimal))
}

func (k Keeper) GetRawAssetPriceFromDenom(ctx sdk.Context, denom string) (sdkmath.LegacyDec, uint64) {
	info, found := k.GetAssetInfo(ctx, denom)
	if !found {
		return sdkmath.LegacyZeroDec(), 0
	}
	price, found := k.GetAssetPrice(ctx, info.Display)
	if !found {
		return sdkmath.LegacyZeroDec(), 0
	}
	return price.Price, info.Decimal
}
