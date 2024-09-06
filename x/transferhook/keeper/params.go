package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/x/transferhook/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.ParamKeyPrefix)
	if b == nil {
		return
	}

	k.Cdc.MustUnmarshal(b, &params)
	return
}

func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.Params) {
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.Cdc.MustMarshal(&params)
	store.Set(types.ParamKeyPrefix, b)
}
