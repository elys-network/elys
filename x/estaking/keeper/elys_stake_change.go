package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/estaking/types"
)

func (k Keeper) SetElysStakeChange(ctx sdk.Context, addr sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetElysStakeChangeKey(addr)
	store.Set(key, addr)
}

func (k Keeper) GetElysStakeChange(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetElysStakeChangeKey(addr)
	return store.Has(key)
}

func (k Keeper) RemoveElysStakeChange(ctx sdk.Context, address sdk.AccAddress) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetElysStakeChangeKey(address)
	store.Delete(key)
}

func (k Keeper) GetAllElysStakeChange(ctx sdk.Context) (list []sdk.AccAddress) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.ElysStakeChangeKeyPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, iterator.Value())
	}

	return
}
