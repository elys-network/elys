package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/burner/types"
)

// SetHistory set a specific history in the store from its index
func (k Keeper) SetHistory(ctx sdk.Context, history types.History) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.HistoryKeyPrefix)
	b := k.cdc.MustMarshal(&history)
	store.Set(types.GetHistoryKey(history.Block), b)
}

// GetHistory returns a history from its index
func (k Keeper) GetHistory(
	ctx sdk.Context,
	block uint64,
) (val types.History, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.HistoryKeyPrefix)

	b := store.Get(types.GetHistoryKey(block))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveHistory removes a history from the store
func (k Keeper) RemoveHistory(
	ctx sdk.Context,
	block uint64,
) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.HistoryKeyPrefix)
	store.Delete(types.GetHistoryKey(block))
}

// GetAllHistory returns all history
func (k Keeper) GetAllHistory(ctx sdk.Context) (list []types.History) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.HistoryKeyPrefix)
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.History
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
