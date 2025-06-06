package keeper

import (
	gomath "math"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) CheckIfWhitelisted(ctx sdk.Context, address sdk.AccAddress) bool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return store.Has(types.GetWhitelistKey(address))
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.GetWhitelistKey(address), address)
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.GetWhitelistKey(address))
}

func (k Keeper) GetWhitelistAddressIterator(ctx sdk.Context) storetypes.Iterator {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	return storetypes.KVStorePrefixIterator(store, types.WhitelistPrefix)
}

func (k Keeper) GetWhitelistedAddress(ctx sdk.Context, pagination *query.PageRequest) ([]sdk.AccAddress, *query.PageResponse, error) {
	var list []sdk.AccAddress
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.WhitelistPrefix)

	if pagination == nil {
		pagination = &query.PageRequest{
			Limit: gomath.MaxUint64 - 1,
		}
	}

	pageRes, err := query.Paginate(prefixStore, pagination, func(key []byte, value []byte) error {
		list = append(list, sdk.AccAddress(value))
		return nil
	})

	return list, pageRes, err
}

func (k Keeper) GetAllWhitelistedAddress(ctx sdk.Context) []sdk.AccAddress {
	var list []sdk.AccAddress
	iterator := k.GetWhitelistAddressIterator(ctx)
	defer func(iterator storetypes.Iterator) {
		err := iterator.Close()
		if err != nil {
			panic(err)
		}
	}(iterator)

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, sdk.AccAddress(iterator.Value()))
	}

	return list
}
