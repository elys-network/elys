package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/perpetual/types"
)

// ResetStore resets all keys in the perpetual module store
func (k Keeper) ResetStore(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	// List of prefixes to clear
	prefixes := [][]byte{
		types.MTPPrefix,
		types.MTPCountPrefix,
		types.OpenMTPCountPrefix,
		types.WhitelistPrefix,
	}

	for _, prefix := range prefixes {
		iter := storetypes.KVStorePrefixIterator(store, prefix)
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			store.Delete(iter.Key())
		}
	}

	return nil
}
