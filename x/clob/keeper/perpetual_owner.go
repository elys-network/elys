package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetualOwner(ctx sdk.Context, owner sdk.AccAddress, marketId uint64) (types.PerpetualOwner, bool) {
	key := types.GetPerpetualOwnerKey(owner, marketId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.PerpetualOwner{}, false
	}

	var v types.PerpetualOwner
	k.cdc.MustUnmarshal(b, &v)
	return v, true
}

func (k Keeper) GetAllSubAccountPerpetuals(ctx sdk.Context) []types.PerpetualOwner {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualOwnerPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualOwner

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualOwner
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetualOwner(ctx sdk.Context, v types.PerpetualOwner) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOwnerKey(v.GetOwnerAccAddress(), v.MarketId)
	b := k.cdc.MustMarshal(&v)
	store.Set(key, b)
}

func (k Keeper) DeletePerpetualOwner(ctx sdk.Context, owner sdk.AccAddress, marketId uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualOwnerKey(owner, marketId)
	store.Delete(key)
}
