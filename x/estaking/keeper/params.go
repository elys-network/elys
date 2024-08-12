package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.ParamsKeyPrefix)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get([]byte(types.LegacyParamsKey))
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

func (k Keeper) DeleteLegacyParams(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	key := []byte(types.LegacyParamsKey)
	store.Delete(key)
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKeyPrefix, b)
}
