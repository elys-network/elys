package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// SetPrice set a specific price in the store from its index
func (k Keeper) SetPrice(ctx sdk.Context, price types.Price) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceKeyPrefix))
	b := k.cdc.MustMarshal(&price)
	store.Set(types.PriceKey(price.Asset), b)
}

// GetPrice returns a price from its index
func (k Keeper) GetPrice(ctx sdk.Context, asset string) (val types.Price, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceKeyPrefix))

	b := store.Get(types.PriceKey(asset))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePrice removes a price from the store
func (k Keeper) RemovePrice(ctx sdk.Context, asset string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.PriceKeyPrefix))
	store.Delete(types.PriceKey(asset))
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
