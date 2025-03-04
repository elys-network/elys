package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetualMarket(ctx sdk.Context, id uint64) (types.PerpetualMarket, error) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketKey(id)
	b := store.Get(key)
	if b == nil {
		return types.PerpetualMarket{}, types.ErrPerpetualMarketNotFound
	}

	var v types.PerpetualMarket
	k.cdc.MustUnmarshal(b, &v)
	return v, nil
}

func (k Keeper) SetPerpetualMarket(ctx sdk.Context, p types.PerpetualMarket) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualMarketKey(p.Id)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}

func (k Keeper) GetAllPerpetualMarket(ctx sdk.Context) []types.PerpetualMarket {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualMarket

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualMarket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) CountAllPerpetualMarket(ctx sdk.Context) uint64 {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()
	count := uint64(0)
	for ; iterator.Valid(); iterator.Next() {
		count++
	}

	return count
}

func (k Keeper) CheckPerpetualMarketAlreadyExists(ctx sdk.Context, baseDenom, quoteDenom string) bool {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualMarketPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualMarket
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		if val.QuoteDenom == quoteDenom && val.BaseDenom == baseDenom {
			return true
		}
	}

	return false
}

// First it checks transient store and then it checks KVstore
func (k Keeper) GetMarketPrice(ctx sdk.Context, id uint64) math.LegacyDec {
	store := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
	key := types.GetMarketPriceKey(id)
	b := store.Get(key)
	if b == nil {
		store = runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		b = store.Get(key)
		if b == nil {
			return math.LegacyZeroDec()
		}
	}

	var v types.MarketPrice
	k.cdc.MustUnmarshal(b, &v)
	return v.LastPrice
}

func (k Keeper) SetMarketPrice(ctx sdk.Context, id uint64, lastPrice math.LegacyDec, permanent bool) {
	store := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
	key := types.GetMarketPriceKey(id)
	v := types.MarketPrice{
		MarketId:  id,
		LastPrice: lastPrice,
	}
	b := k.cdc.MustMarshal(&v)
	store.Set(key, b)

	if permanent {
		store = runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		store.Set(key, b)
	}
}

func (k Keeper) GetAllMarketPrice(ctx sdk.Context) []types.MarketPrice {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.MarketPricePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.MarketPrice

	for ; iterator.Valid(); iterator.Next() {
		var val types.MarketPrice
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}
