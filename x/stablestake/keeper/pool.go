package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) GetPool(ctx sdk.Context, id uint64) types.Pool {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolKey(id)
	bz := store.Get(key)
	if bz == nil {
		return types.Pool{
			Id:               id,
			TotalLiabilities: sdk.Coins{},
		}
	}
	var val types.Pool
	k.cdc.MustUnmarshal(bz, &val)
	return val
}

func (k Keeper) SetPool(ctx sdk.Context, pool types.Pool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPoolKey(pool.Id)
	bz := k.cdc.MustMarshal(&pool)
	store.Set(key, bz)
}

func (k Keeper) AddPoolLiabilities(ctx sdk.Context, id uint64, coin sdk.Coin) {
	pool := k.GetPool(ctx, id)
	pool.AddLiabilities(coin)
	k.SetPool(ctx, pool)
}
