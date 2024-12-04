package keeper

import (
	"cosmossdk.io/math"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/estaking/types"
)

// GetParams get all parameters as types.Params
func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsKeyPrefix)
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
	store.Set(types.ParamsKeyPrefix, b)
}

func (k Keeper) GetLegacyParams(ctx sdk.Context) (params types.LegacyParams) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get(types.ParamsKeyPrefix)
	if b == nil {
		return
	}

	k.cdc.MustUnmarshal(b, &params)
	return
}

func (k Keeper) SetDexRewardsStakers(ctx sdk.Context, amount math.LegacyDec) {
	params := k.GetParams(ctx)
	// Set block number and total dex rewards given
	params.DexRewardsStakers.NumBlocks = 1 // Don't know why this is 1
	params.DexRewardsStakers.Amount = amount
	k.SetParams(ctx, params)
}
