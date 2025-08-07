package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
)

func (k Keeper) SetFallbackCounter(ctx sdk.Context, v types.FallbackCounter) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.FallbackCounterPrefix, k.cdc.MustMarshal(&v))
}

func (k Keeper) GetFallbackCounter(ctx sdk.Context) types.FallbackCounter {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	bz := store.Get(types.FallbackCounterPrefix)
	if bz == nil {
		return types.FallbackCounter{
			Counter: 0,
			NextKey: nil,
		}
	} else {
		var v types.FallbackCounter
		k.cdc.MustUnmarshal(bz, &v)
		return v
	}
}

func (k Keeper) DeleteLegacyFallbackOffset(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete(types.LegacyFallbackOffsetKeyPrefix)
}
