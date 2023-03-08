package keeper

import (
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tokenomics/types"
)

// SetGenesisInflation set genesisInflation in the store
func (k Keeper) SetGenesisInflation(ctx sdk.Context, genesisInflation types.GenesisInflation) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GenesisInflationKey))
	b := k.cdc.MustMarshal(&genesisInflation)
	store.Set([]byte{0}, b)
}

// GetGenesisInflation returns genesisInflation
func (k Keeper) GetGenesisInflation(ctx sdk.Context) (val types.GenesisInflation, found bool) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GenesisInflationKey))

	b := store.Get([]byte{0})
	if b == nil {
		return val, false
	}

	k.cdc.MustUnmarshal(b, &val)
	return val, true
}

// RemoveGenesisInflation removes genesisInflation from the store
func (k Keeper) RemoveGenesisInflation(ctx sdk.Context) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.KeyPrefix(types.GenesisInflationKey))
	store.Delete([]byte{0})
}
