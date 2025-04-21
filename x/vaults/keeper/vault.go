package keeper

import (
	storetypes "cosmossdk.io/store/types"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/vaults/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetVault(ctx sdk.Context, id uint64) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.GetVaultKey(id))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetVault(ctx sdk.Context, vault types.Vault) error {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&vault)
	store.Set(types.GetVaultKey(vault.Id), b)
	return nil
}

func (k Keeper) GetLatestVault(ctx sdk.Context) (val types.Vault, found bool) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	iterator := storetypes.KVStoreReversePrefixIterator(store, types.VaultKeyPrefix)
	defer iterator.Close()

	if !iterator.Valid() {
		return val, false
	}

	k.cdc.MustUnmarshal(iterator.Value(), &val)
	return val, true
}

// GetNextVaultId returns the next vault id.
func (k Keeper) GetNextVaultId(ctx sdk.Context) uint64 {
	latestVault, found := k.GetLatestVault(ctx)
	if !found {
		return 1
	}
	return latestVault.Id + 1
}
