package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/commitment/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// GetLegacyParams get all parameters as types.Params
func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsKey)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

// SetParams set the params
func (k Keeper) SetParams(ctx sdk.Context, params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	b := k.cdc.MustMarshal(&params)
	store.Set(types.ParamsKey, b)
}

// GetVestingDenom returns the vesting denom for the given base denom
func (k Keeper) GetVestingInfo(ctx sdk.Context, baseDenom string) (*types.VestingInfo, int) {
	params := k.GetParams(ctx)

	for i, vestingInfo := range params.VestingInfos {
		if vestingInfo.BaseDenom == baseDenom {
			return &vestingInfo, i
		}
	}

	return nil, 0
}
