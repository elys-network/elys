package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/osmomath"
	"github.com/elys-network/elys/x/oracle/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"github.com/osmosis-labs/osmosis/v13/osmomath"
)

// SetPrice set a specific price in the store from its index
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&price)
	store.Set(types.PriceKey(price.Asset, price.Source, price.Timestamp), b)
}

// GetPrice returns a price from its index
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

// RemovePrice removes a price from the store
func (k Keeper) RemovePrice(ctx sdk.Context, asset, source string, timestamp uint64) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PriceKey(asset, source, timestamp))
}

// GetAllPrice returns all price
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

func (k Keeper) GetAssetPrice(ctx sdk.Context, asset string) (types.Price, bool) {
	// try out band source
	val, found := k.GetLatestPriceFromAssetAndSource(ctx, asset, types.BAND)
	if found {
		return val, true
	}

	// try out elys source
	val, found = k.GetLatestPriceFromAssetAndSource(ctx, asset, types.ELYS)
	if found {
		return val, true
	}

	// find from any source if band source does not exist
	return k.GetLatestPriceFromAnySource(ctx, asset)
}

func (k Keeper) GetAssetPriceFromDenom(ctx sdk.Context, denom string) sdk.Dec {
	info, found := k.GetAssetInfo(ctx, denom)
	if !found {
		return sdk.ZeroDec()
	}
	price, found := k.GetAssetPrice(ctx, info.Display)
	if !found {
		return sdk.ZeroDec()
	}
	ten := osmomath.NewBigDec(10)
	return price.Price.Quo(osmomath.Pow(ten, info.Decimals))
}
