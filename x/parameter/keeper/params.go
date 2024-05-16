package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/parameter/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get([]byte(types.ParamsKey))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set([]byte(types.ParamsKey), b)
}

// GetParams get all parameters as types.Params
func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get([]byte(types.ParamsKey))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}
