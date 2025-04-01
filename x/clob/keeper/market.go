package keeper

import (
	"cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
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
func (k Keeper) GetLastMarketPrice(ctx sdk.Context, id uint64) math.Dec {
	tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
	key := types.GetLastMarketPriceKey(id)
	b := tStore.Get(key)
	if b == nil {
		store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		b = store.Get(key)
		if b == nil {
			return utils.ZeroDec
		}
		tStore.Set(key, b)
	}

	var v types.LastMarketPrice
	k.cdc.MustUnmarshal(b, &v)
	return v.LastPrice
}

func (k Keeper) SetLastMarketPrice(ctx sdk.Context, id uint64, lastPrice math.Dec, set bool) {
	tStore := runtime.KVStoreAdapter(k.transientStoreService.OpenTransientStore(ctx))
	key := types.GetLastMarketPriceKey(id)
	v := types.LastMarketPrice{
		MarketId:  id,
		LastPrice: lastPrice,
	}
	b := k.cdc.MustMarshal(&v)
	tStore.Set(key, b)

	if set {
		store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
		store.Set(key, b)
	}
}

func (k Keeper) GetAllLastMarketPrice(ctx sdk.Context) []types.LastMarketPrice {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.LastMarketPricePrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.LastMarketPrice

	for ; iterator.Valid(); iterator.Next() {
		var val types.LastMarketPrice
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}
