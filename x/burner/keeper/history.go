package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/burner/types"
)

// SetHistory set a specific history in the store from its index
func (k Keeper) SetHistory(ctx sdk.Context, history types.History) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.HistoryKeyPrefix))
	b := k.cdc.MustMarshal(&history)
	store.Set(types.HistoryKey(
		history.Timestamp,
		history.Denom,
	), b)
}

// GetHistory returns a history from its index
func (k Keeper) GetHistory(
	ctx sdk.Context,
	timestamp string,
	denom string,
) (val types.History, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.HistoryKeyPrefix))

	b := store.Get(types.HistoryKey(
		timestamp,
		denom,
	))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveHistory removes a history from the store
func (k Keeper) RemoveHistory(
	ctx sdk.Context,
	timestamp string,
	denom string,
) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.HistoryKeyPrefix))
	store.Delete(types.HistoryKey(
		timestamp,
		denom,
	))
}

// GetAllHistory returns all history
func (k Keeper) GetAllHistory(ctx sdk.Context) (list []types.History) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.HistoryKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.History
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
