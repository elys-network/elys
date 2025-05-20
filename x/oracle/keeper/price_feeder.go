package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/oracle/types"
)

// SetPriceFeeder set a specific priceFeeder in the store from its index
func (k Keeper) SetPriceFeeder(ctx sdk.Context, priceFeeder types.PriceFeeder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPriceFeederKey(priceFeeder.GetFeederAccount())
	b := k.cdc.MustMarshal(&priceFeeder)
	store.Set(key, b)
}

// GetPriceFeeder returns a priceFeeder from its index
func (k Keeper) GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val types.PriceFeeder, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPriceFeederKey(feeder)

	b := store.Get(key)
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePriceFeeder removes a priceFeeder from the store
func (k Keeper) RemovePriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPriceFeederKey(feeder)
	store.Delete(key)
}

// GetAllPriceFeeder returns all priceFeeder
func (k Keeper) GetAllPriceFeeder(ctx sdk.Context) (list []types.PriceFeeder) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.PriceFeederPrefixKey)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PriceFeeder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAllLegacyPriceFeeder(ctx sdk.Context) (list []types.PriceFeeder) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.LegacyPriceFeederKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PriceFeeder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) RemoveLegacyPriceFeeder(ctx sdk.Context, feeder string) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.LegacyPriceFeederKeyPrefix))
	store.Delete(types.LegacyPriceFeederKey(feeder))
}
