package keeper

import (
	"errors"

	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetOrderOwner(ctx sdk.Context, owner sdk.AccAddress, subAccountId uint64, orderKey types.OrderKey) (types.PerpetualOrderOwner, error) {
	key := types.GetOrderOwnerKey(owner, subAccountId, orderKey)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.PerpetualOrderOwner{}, errors.New("order owner not found")
	}

	var val types.PerpetualOrderOwner
	k.cdc.MustUnmarshal(b, &val)
	return val, nil
}

func (k Keeper) GetAllOrdersForOwner(ctx sdk.Context, owner sdk.AccAddress) []types.PerpetualOrderOwner {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOrderOwnerPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualOrderOwner

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrderOwner
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) GetAllOrderOwners(ctx sdk.Context) []types.PerpetualOrderOwner {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOrderOwnerPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualOrderOwner

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOrderOwner
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetOrderOwner(ctx sdk.Context, v types.PerpetualOrderOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetOrderOwnerKey(v.GetOwnerAccAddress(), v.GetSubAccountId(), v.GetOrderKey())
	b := k.cdc.MustMarshal(&v)
	store.Set(key, b)
}

func (k Keeper) DeleteOrderOwner(ctx sdk.Context, v types.PerpetualOrderOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetOrderOwnerKey(v.GetOwnerAccAddress(), v.GetSubAccountId(), v.GetOrderKey())
	store.Delete(key)
}
