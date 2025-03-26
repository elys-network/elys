package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetAmmPool(ctx sdk.Context, id uint64) types.AmmPool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetAmmPoolKey(id)
	bz := store.Get(key)
	if bz == nil {
		return types.AmmPool{
			Id:               id,
			TotalLiabilities: sdk.Coins{},
		}
	}
	var val types.AmmPool
	k.cdc.MustUnmarshal(bz, &val)
	return val
}

func (k Keeper) GetAllAmmPools(ctx sdk.Context) []types.AmmPool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStorePrefixIterator(store, types.AmmPoolKeyPrefix)

	defer iterator.Close()
	var vals []types.AmmPool
	for ; iterator.Valid(); iterator.Next() {
		var val types.AmmPool
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		vals = append(vals, val)
	}
	return vals
}

func (k Keeper) SetAmmPool(ctx sdk.Context, pool types.AmmPool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetAmmPoolKey(pool.Id)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
}

func (k Keeper) AddPoolLiabilities(ctx sdk.Context, id uint64, coin sdk.Coin) {
	pool := k.GetAmmPool(ctx, id)
	pool.AddLiabilities(coin)
	k.SetAmmPool(ctx, pool)
}

func (k Keeper) SubtractPoolLiabilities(ctx sdk.Context, id uint64, coin sdk.Coin) {
	pool := k.GetAmmPool(ctx, id)
	pool.SubLiabilities(coin)
	k.SetAmmPool(ctx, pool)
}
