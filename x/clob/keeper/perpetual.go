package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) GetPerpetualOwner(ctx sdk.Context, subAccountId uint64, owner sdk.AccAddress, marketId uint64) (types.PerpetualOwner, bool) {
	key := types.GetPerpetualOwnerKey(subAccountId, owner, marketId)
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
	key := types.GetPerpetualOwnerKey(v.SubAccountId, v.GetOwnerAccAddress(), v.MarketId)
	b := k.cdc.MustMarshal(&v)
	store.Set(key, b)
}

func (k Keeper) GetPerpetual(ctx sdk.Context, marketId, id uint64) (types.Perpetual, error) {
	key := types.GetPerpetualKey(marketId, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.Perpetual{}, types.ErrPerpetualNotFound
	}

	var v types.Perpetual
	k.cdc.MustUnmarshal(b, &v)
	return v, nil
}

func (k Keeper) GetAllPerpetuals(ctx sdk.Context) []types.Perpetual {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.Perpetual

	for ; iterator.Valid(); iterator.Next() {
		var val types.Perpetual
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetual(ctx sdk.Context, p types.Perpetual) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualKey(p.MarketId, p.Id)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}

func (k Keeper) DeletePerpetual(ctx sdk.Context, p types.Perpetual) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualKey(p.MarketId, p.Id)
	store.Delete(key)
}

func (k Keeper) GetAndUpdatePerpetualCounter(ctx sdk.Context, marketId uint64) uint64 {
	key := types.GetPerpetualCounterKey(marketId)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		v := types.PerpetualCounter{
			MarketId: marketId,
			Counter:  2,
		}
		b = k.cdc.MustMarshal(&v)
		store.Set(key, b)
		return 1
	}

	var v types.PerpetualCounter
	k.cdc.MustUnmarshal(b, &v)
	result := v.Counter
	v.Counter = v.Counter + 1
	b = k.cdc.MustMarshal(&v)
	store.Set(key, b)
	return result
}

func (k Keeper) GetAllPerpetualCounters(ctx sdk.Context) []types.Perpetual {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualCounterPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.Perpetual

	for ; iterator.Valid(); iterator.Next() {
		var val types.Perpetual
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) setPerpetualCounter(ctx sdk.Context, p types.PerpetualCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualCounterKey(p.MarketId)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}
