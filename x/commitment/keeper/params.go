package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

// GetParams returns the current parameters of the Commitment module
func (k Keeper) GetLegacyParams(ctx sdk.Context) types.LegacyParams {
	var params types.LegacyParams
	k.paramstore.GetParamSet(ctx, &params)
	return params
}

// SetParams set the params
func (k Keeper) SetLegacyParams(ctx sdk.Context, params types.LegacyParams) {
	k.paramstore.SetParamSet(ctx, &params)
}

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
