package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/accountedpool/types"
)

// SetAccountedPool set a specific accountedPool in the store from its index
func (k Keeper) SetAccountedPool(ctx sdk.Context, accountedPool types.AccountedPool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
	b := k.cdc.MustMarshal(&accountedPool)
	store.Set(types.AccountedPoolKey(accountedPool.PoolId), b)
}

// GetAccountedPool returns a accountedPool from its index
func (k Keeper) GetAccountedPool(ctx sdk.Context, PoolId uint64) (val types.AccountedPool, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))

	b := store.Get(types.AccountedPoolKey(PoolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveAccountedPool removes a accountedPool from the store
func (k Keeper) RemoveAccountedPool(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
	store.Delete(types.AccountedPoolKey(poolId))
}

// GetAllAccountedPool returns all accountedPool
func (k Keeper) GetAllAccountedPool(ctx sdk.Context) (list []types.AccountedPool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.AccountedPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetAllLegacyAccountedPool returns all legacyAccountedPool
func (k Keeper) GetAllLegacyAccountedPool(ctx sdk.Context) (list []types.LegacyAccountedPool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyAccountedPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// PoolExists checks if a pool with the given poolId exists in the list of pools
func (k Keeper) PoolExists(ctx sdk.Context, poolId uint64) bool {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.AccountedPoolKeyPrefix))
	b := store.Get(types.AccountedPoolKey(poolId))
	return b != nil
}
