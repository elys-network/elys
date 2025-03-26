package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, index uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(index))
}

// GetAllPools returns all pool
func (k Keeper) GetAllPools(ctx sdk.Context) (list []types.Pool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// SetPool set a specific pool in the store from its index
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKey(pool.AmmPoolId), b)
}

func (k Keeper) DeletePool(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(poolId))
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(
		poolId,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

func (k Keeper) SetLeveragedAmount(ctx sdk.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		ammPool, found := k.amm.GetPool(ctx, pool.AmmPoolId)
		if !found {
			continue
		}
		for _, asset := range ammPool.PoolAssets {
			pool.AssetLeverageAmounts = append(pool.AssetLeverageAmounts, &types.AssetLeverageAmount{
				Denom:           asset.Token.Denom,
				LeveragedAmount: sdkmath.ZeroInt(),
			})
		}
		k.SetPool(ctx, pool)
	}

	iterator := k.GetPositionIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var position types.Position
		k.cdc.MustUnmarshal(iterator.Value(), &position)
		pool, found := k.GetPool(ctx, position.AmmPoolId)
		if !found {
			continue
		}
		pool.UpdateAssetLeveragedAmount(ctx, position.Collateral.Denom, position.LeveragedLpAmount, true)
		k.SetPool(ctx, pool)
	}
}
