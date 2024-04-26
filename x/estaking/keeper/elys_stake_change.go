package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

func (k Keeper) SetElysStakeChange(ctx sdk.Context, addr sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakeChangeKeyPrefix))
	store.Set([]byte(addr), addr)
}

func (k Keeper) GetElysStakeChange(ctx sdk.Context, addr sdk.AccAddress) (found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakeChangeKeyPrefix))

	return store.Has([]byte(addr))
}

func (k Keeper) RemoveElysStakeChange(ctx sdk.Context, address sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakeChangeKeyPrefix))
	store.Delete([]byte(address))
}

func (k Keeper) GetAllElysStakeChange(ctx sdk.Context) (list []sdk.AccAddress) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.ElysStakeChangeKeyPrefix))
	iterator := sdk.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		list = append(list, sdk.AccAddress(iterator.Value()))
	}

	return
}
