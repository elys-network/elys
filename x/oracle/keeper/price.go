package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/osmosis-labs/osmosis/osmomath"
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

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (osmomath.BigDec, bool) {
	// try out elys source
	val, found := k.GetLatestPriceFromAssetAndSource(ctx, asset, types.ELYS)
	if found {
		return osmomath.BigDecFromDec(val.Price), true
	}

	// try out band source
	val, found = k.GetLatestPriceFromAssetAndSource(ctx, asset, types.BAND)
	if found {
		return osmomath.BigDecFromDec(val.Price), true
	}

	// find from any source if band source does not exist
	price, found := k.GetLatestPriceFromAnySource(ctx, asset)
	if found {
		return osmomath.BigDecFromDec(price.Price), true
	}
	return osmomath.BigDec{}, false
}

func (k Keeper) GetDenomPrice(ctx sdk.Context, denom string) osmomath.BigDec {
	info, found := k.GetAssetInfo(ctx, denom)
	if !found {
		return osmomath.ZeroBigDec()
	}
	price, found := k.GetAssetPrice(ctx, info.Display)
	if !found {
		return osmomath.ZeroBigDec()
	}
	if info.Decimal <= 18 {
		return price.QuoInt64(utils.Pow10Int64(info.Decimal))
	}
	return price.Quo(utils.Pow10(info.Decimal))
}
