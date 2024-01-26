package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// ResetStore resets all keys in the perpetual module store
func (k Keeper) ResetStore(ctx sdk.Context) error {
	store := ctx.KVStore(k.storeKey)

	// List of prefixes to clear
	prefixes := [][]byte{
		types.MTPPrefix,
		types.MTPCountPrefix,
		types.OpenMTPCountPrefix,
		types.WhitelistPrefix,
	}

	for _, prefix := range prefixes {
		iter := sdk.KVStorePrefixIterator(store, prefix)
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			store.Delete(iter.Key())
		}
	}

	return nil
}
