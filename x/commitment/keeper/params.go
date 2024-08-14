package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := ctx.KVStore(k.storeKey)

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// remove after migration
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
	store.Delete([]byte(types.LegacyParamsKey))
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := ctx.KVStore(k.storeKey)
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, b)
}

// GetVestingDenom returns the vesting denom for the given base denom
func (k Keeper) GetVestingInfo(ctx sdk.Context, baseDenom string) (*types.VestingInfo, int) {
	params := k.GetParams(ctx)

	for i, vestingInfo := range params.VestingInfos {
		if vestingInfo.BaseDenom == baseDenom {
			return vestingInfo, i
		}
	}

	return nil, 0
}
