package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

func (k Keeper) SetElysStakeChange(ctx sdk.Context, addr sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetElysStakeChangeKey(addr)
	store.Set(key, addr)
}

func (k Keeper) GetElysStakeChange(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetElysStakeChangeKey(addr)
	return store.Has(key)
}

func (k Keeper) RemoveElysStakeChange(ctx sdk.Context, address sdk.AccAddress) {
	store := ctx.KVStore(k.storeKey)
	key := types.GetElysStakeChangeKey(address)
	store.Delete(key)
}

func (k Keeper) DeleteLegacyElysStakeChange(ctx sdk.Context, address sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyElysStakeChangeKeyPrefix))
	store.Delete([]byte(address))
}

func (k Keeper) GetAllElysStakeChange(ctx sdk.Context) (list []sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.ElysStakeChangeKeyPrefix)
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Value())
	}

	return
}

func (k Keeper) GetAllLegacyElysStakeChange(ctx sdk.Context) (list []sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.LegacyKeyPrefix(types.LegacyElysStakeChangeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, sdk.AccAddress(iterator.Value()))
	}

	return
}
