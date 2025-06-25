package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/tokenomics/types"
)

// SetGenesisInflation set genesisInflation in the store
func (k Keeper) SetGenesisInflation(ctx sdk.Context, genesisInflation types.GenesisInflation) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&genesisInflation)
	store.Set([]byte(types.GenesisInflationKey), b)
}

// GetGenesisInflation returns genesisInflation
func (k Keeper) GetGenesisInflation(ctx sdk.Context) (val types.GenesisInflation, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get([]byte(types.GenesisInflationKey))
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGenesisInflation removes genesisInflation from the store
func (k Keeper) RemoveGenesisInflation(ctx sdk.Context) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Delete([]byte(types.GenesisInflationKey))
}
