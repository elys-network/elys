package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

func (k Keeper) GetPerpetualADL(ctx sdk.Context, marketId, id uint64) (types.PerpetualADL, bool) {
	key := types.GetPerpetualADLKey(marketId, id)
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(key)
	if b == nil {
		return types.PerpetualADL{}, false
	}

	var v types.PerpetualADL
	k.cdc.MustUnmarshal(b, &v)
	return v, true
}

func (k Keeper) GetAllPerpetualADLs(ctx sdk.Context) []types.PerpetualADL {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.PerpetualADLPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	var list []types.PerpetualADL

	for ; iterator.Valid(); iterator.Next() {
		var val types.PerpetualADL
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return list
}

func (k Keeper) SetPerpetualADL(ctx sdk.Context, p types.PerpetualADL) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualADLKey(p.MarketId, p.Id)
	b := k.cdc.MustMarshal(&p)
	store.Set(key, b)
}

func (k Keeper) DeletePerpetualADL(ctx sdk.Context, marketId, id uint64) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := types.GetPerpetualADLKey(marketId, id)
	store.Delete(key)
}
