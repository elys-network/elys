package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/stablestake/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamKeyPrefix)
	if b == nil {
		// TODO fix this
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamKeyPrefix, b)
}

// GetLegacyDepositDenom deprecated
func (k Keeper) GetLegacyDepositDenom(ctx sdk.Context) string {
	params := k.GetParams(ctx)
	entry, found := k.assetProfileKeeper.GetEntry(ctx, params.LegacyDepositDenom)
	if !found {
		return params.LegacyDepositDenom
	}
	return entry.Denom
}
