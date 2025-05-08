package keeper

import (
	sdkmath "cosmossdk.io/math"
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// SetPool set a specific pool in the store from its index
func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	err := pool.Validate()
	if err != nil {
		panic(err)
	}
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	key := types.PoolKey(pool.PoolId)
	store.Set(key, b)
}

func (k Keeper) SetLegacyPool(ctx sdk.Context, pool types.LegacyPool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	b := k.cdc.MustMarshal(&pool)
	store.Set(types.PoolKey(pool.PoolId), b)
}

// GetPool returns a pool from its index
func (k Keeper) GetPool(ctx sdk.Context, poolId uint64) (val types.Pool, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))

	b := store.Get(types.PoolKey(poolId))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemovePool removes a pool from the store
func (k Keeper) RemovePool(ctx sdk.Context, poolId uint64) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	store.Delete(types.PoolKey(poolId))
}

// GetAllPool returns all pool
func (k Keeper) GetAllPool(ctx sdk.Context) (list []types.Pool) {
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

// GetAllLegacyPool returns all legacy pool
func (k Keeper) GetAllLegacyPool(ctx sdk.Context) (list []types.LegacyPool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.LegacyPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

// GetLatestPool retrieves the latest pool item from the list of pools
func (k Keeper) GetLatestPool(ctx sdk.Context) (val types.Pool, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := storetypes.KVStoreReversePrefixIterator(store, []byte{})
	defer iterator.Close()

	if !iterator.Valid() {
		return val, false
	}

	k.cdc.MustUnmarshal(iterator.Value(), &val)
	return val, true
}

// GetNextPoolId returns the next pool id.
func (k Keeper) GetNextPoolId(ctx sdk.Context) uint64 {
	latestPool, found := k.GetLatestPool(ctx)
	if !found {
		return 1
	}
	return latestPool.PoolId + 1
}

// PoolExists checks if a pool with the given poolId exists in the list of pools
func (k Keeper) PoolExists(ctx sdk.Context, poolId uint64) bool {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	b := store.Get(types.PoolKey(poolId))
	return b != nil
}

// GetBestPoolWithDenoms returns the first highest TVL pool id that contains all specified denominations
func (k Keeper) GetBestPoolWithDenoms(ctx sdk.Context, denoms []string, usesOracle bool) (pool types.Pool, found bool) {
	// Get all pools
	pools := k.GetAllPool(ctx)

	maxTvl := sdkmath.LegacyNewDec(-1)
	bestPool := types.Pool{}
	for _, p := range pools {
		// If usesOracle is false, function filters in all pools.
		if usesOracle && !p.PoolParams.UseOracle {
			continue
		}
		// If the number of assets in the pool is less than the number of denoms, skip
		if len(p.PoolAssets) < len(denoms) {
			continue
		}

		allDenomsFound := true

		// Check that all denoms are in the pool's assets
		for _, denom := range denoms {
			denomFound := false
			for _, asset := range p.PoolAssets {
				if denom == asset.Token.Denom {
					denomFound = true
					break
				}
			}

			// If any denom is not found, mark allDenomsFound as false and break
			if !denomFound {
				allDenomsFound = false
				break
			}
		}

		poolTvl, err := p.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			poolTvl = osmomath.ZeroBigDec()
		}

		// If all denoms are found in this pool, return the pool id
		if allDenomsFound && maxTvl.LT(poolTvl.Dec()) {
			maxTvl = poolTvl.Dec()
			bestPool = p
		}
	}

	return bestPool, !maxTvl.IsNegative()
}

// IterateLiquidityPools iterates over all LiquidityPools and performs a
// callback.
func (k Keeper) IterateLiquidityPools(ctx sdk.Context, handlerFn func(pool types.Pool) (stop bool)) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.PoolKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var pool types.Pool
		k.cdc.MustUnmarshal(iterator.Value(), &pool)

		if handlerFn(pool) {
			break
		}
	}
}

// GetPoolWithAccountedBalance Gets the pool snapshot and updates the pool balance with accounted pool balance
func (k Keeper) GetPoolWithAccountedBalance(ctx sdk.Context, poolId uint64) (val types.SnapshotPool) {
	snapshot, found := k.GetPool(ctx, poolId)
	if !found {
		panic(fmt.Sprintf("pool %d not found", poolId))
	}
	poolAssets := []types.PoolAsset{}
	// Update the pool snapshot with accounted pool balance
	for _, asset := range snapshot.PoolAssets {
		accAmount := k.accountedPoolKeeper.GetAccountedBalance(ctx, poolId, asset.Token.Denom)
		if accAmount.IsPositive() {
			asset.Token.Amount = accAmount
		}
		poolAssets = append(poolAssets, asset)
	}
	snapshot.PoolAssets = poolAssets

	return types.SnapshotPool{Pool: snapshot}
}

// AddToPoolBalance Used in perpetual balance changes
func (k Keeper) AddToPoolBalanceAndUpdateLiquidity(ctx sdk.Context, pool *types.Pool, addShares sdkmath.Int, coins sdk.Coins) error {
	err := pool.IncreaseLiquidity(addShares, coins)
	if err != nil {
		return err
	}
	k.SetPool(ctx, *pool)
	return k.RecordTotalLiquidityIncrease(ctx, coins)
}

// RemoveFromPoolBalance Used in perpetual balance changes
func (k Keeper) RemoveFromPoolBalanceAndUpdateLiquidity(ctx sdk.Context, pool *types.Pool, removeShares sdkmath.Int, coins sdk.Coins) error {
	err := pool.DecreaseLiquidity(removeShares, coins)
	if err != nil {
		return err
	}
	k.SetPool(ctx, *pool)
	return k.RecordTotalLiquidityDecrease(ctx, coins)
}

// For migration only, fixes for amm mismatch has been added, to  verify if the fix is working we
// will match the balances and do operations after
func (k Keeper) MatchAmmBalances(ctx sdk.Context) error {
	// Match pool balances and assets structure balances
	pools := k.GetAllPool(ctx)
	for _, pool := range pools {
		if !pool.PoolParams.UseOracle {
			balances := k.bankKeeper.GetAllBalances(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()))
			for _, asset := range pool.PoolAssets {
				for _, balance := range balances {
					if asset.Token.Denom == balance.Denom {
						if asset.Token.Amount.GT(balance.Amount) {
							k.bankKeeper.MintCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Token.Denom, asset.Token.Amount.Sub(balance.Amount))))
							k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(pool.GetAddress()), sdk.NewCoins(sdk.NewCoin(asset.Token.Denom, asset.Token.Amount.Sub(balance.Amount))))
						}

						if asset.Token.Amount.LT(balance.Amount) {
							k.bankKeeper.SendCoinsFromAccountToModule(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()), types.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Token.Denom, balance.Amount.Sub(asset.Token.Amount))))
							k.bankKeeper.BurnCoins(ctx, types.ModuleName, sdk.NewCoins(sdk.NewCoin(asset.Token.Denom, balance.Amount.Sub(asset.Token.Amount))))
						}
					}
				}
			}
		}
	}
	return nil
}
