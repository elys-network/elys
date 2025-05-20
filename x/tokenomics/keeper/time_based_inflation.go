package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v4/x/tokenomics/types"
)

// SetTimeBasedInflation set a specific timeBasedInflation in the store from its index
func (k Keeper) SetTimeBasedInflation(ctx sdk.Context, timeBasedInflation types.TimeBasedInflation) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.TimeBasedInflationKeyPrefix))
	b := k.cdc.MustMarshal(&timeBasedInflation)
	store.Set(types.TimeBasedInflationKey(
		timeBasedInflation.StartBlockHeight,
		timeBasedInflation.EndBlockHeight,
	), b)
}

// GetTimeBasedInflation returns a timeBasedInflation from its index
func (k Keeper) GetTimeBasedInflation(
	ctx sdk.Context,
	startBlockHeight uint64,
	endBlockHeight uint64,
) (val types.TimeBasedInflation, found bool) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.TimeBasedInflationKeyPrefix))

	b := store.Get(types.TimeBasedInflationKey(startBlockHeight, endBlockHeight))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveTimeBasedInflation removes a timeBasedInflation from the store
func (k Keeper) RemoveTimeBasedInflation(
	ctx sdk.Context,
	startBlockHeight uint64,
	endBlockHeight uint64,
) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.TimeBasedInflationKeyPrefix))
	store.Delete(types.TimeBasedInflationKey(startBlockHeight, endBlockHeight))
}

// GetAllTimeBasedInflation returns all timeBasedInflation
func (k Keeper) GetAllTimeBasedInflation(ctx sdk.Context) (list []types.TimeBasedInflation) {
	store := prefix.NewStore(runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)), types.KeyPrefix(types.TimeBasedInflationKeyPrefix))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.TimeBasedInflation
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}
