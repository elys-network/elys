package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// SetPriceFeeder set a specific priceFeeder in the store from its index
func (k Keeper) SetPriceFeeder(ctx sdk.Context, priceFeeder types.PriceFeeder) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceFeederKey(priceFeeder.GetFeederAccount())
	b := k.cdc.MustMarshal(&priceFeeder)
	store.Set(key, b)
}

// GetPriceFeeder returns a priceFeeder from its index
func (k Keeper) GetPriceFeeder(ctx sdk.Context, feeder sdk.AccAddress) (val types.PriceFeeder, found bool) {
	store := ctx.KVStore(k.storeKey)
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
	store := ctx.KVStore(k.storeKey)
	key := types.GetPriceFeederKey(feeder)
	store.Delete(key)
}

// GetAllPriceFeeder returns all priceFeeder
func (k Keeper) GetAllPriceFeeder(ctx sdk.Context) (list []types.PriceFeeder) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.PriceFeederPrefixKey)

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PriceFeeder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetAllLegacyPriceFeeder(ctx sdk.Context) (list []types.PriceFeeder) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyPriceFeederKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.PriceFeeder
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) RemoveLegacyPriceFeeder(ctx sdk.Context, feeder string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.LegacyPriceFeederKeyPrefix))
	store.Delete(types.LegacyPriceFeederKey(feeder))
}
