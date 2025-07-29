package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

// ResetStore resets all keys in the perpetual module store
func (k Keeper) ResetStore(ctx sdk.Context) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	// List of prefixes to clear
	prefixes := [][]byte{
		types.MTPPrefix,
		types.LegacyMTPCountPrefix,
		types.LegacyOpenMTPCountPrefix,
		types.WhitelistPrefix,
		types.PerpetualCounterPrefix,
	}

	for _, prefix := range prefixes {
		iter := storetypes.KVStorePrefixIterator(store, prefix)
		defer iter.Close()

		for ; iter.Valid(); iter.Next() {
			store.Delete(iter.Key())
		}
	}

	params := k.GetParams(ctx)
	allPools := k.GetAllPools(ctx)
	for _, pool := range allPools {
		ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
		if err != nil {
			return err
		}
		resetPool := types.NewPool(ammPool, params.LeverageMax)
		k.SetPool(ctx, resetPool)
	}
	return nil
}
