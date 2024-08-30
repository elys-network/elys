package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/perpetual/types"
	gomath "math"
)

func (k Keeper) CheckIfWhitelisted(ctx sdk.Context, address sdk.AccAddress) bool {
	store := ctx.KVStore(k.storeKey)
	return store.Has(types.GetWhitelistKey(address))
}

func (k Keeper) WhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetWhitelistKey(address), address)
}

func (k Keeper) DewhitelistAddress(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetWhitelistKey(address))
}

func (k Keeper) GetWhitelistAddressIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, types.WhitelistPrefix)
}

func (k Keeper) GetWhitelistedAddress(ctx sdk.Context, pagination *query.PageRequest) ([]sdk.AccAddress, *query.PageResponse, error) {
	var list []sdk.AccAddress
	store := ctx.KVStore(k.storeKey)
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
	defer func(iterator sdk.Iterator) {
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

func (k Keeper) V6_MigrateWhitelistedAddress(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	iterator := sdk.KVStorePrefixIterator(store, types.WhitelistPrefix)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		address := (string)(iterator.Value())
		accAddress := sdk.MustAccAddressFromBech32(address)
		k.removeLegacyWhitelistAddress(ctx, address)
		k.WhitelistAddress(ctx, accAddress)
	}

	return
}

func (k Keeper) removeLegacyWhitelistAddress(ctx sdk.Context, address string) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.GetLegacyWhitelistKey(address))
}
