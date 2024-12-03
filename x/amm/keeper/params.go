package keeper

import (
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) GetParams(ctx sdk.Context) (params types.Params) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))

	b := store.Get([]byte(types.ParamsKey))
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
	store.Set([]byte(types.ParamsKey), b)
}

func (k Keeper) CheckBaseAssetExist(ctx sdk.Context, denom string) bool {

	params := k.GetParams(ctx)

	// We need to do this step because when initializing chain, usdc denom will be unknown until ibc is set up.
	// Then adding usdc denom through gov proposal will take time, and we won't be able to open a pool until proposal gets executed
	if len(params.BaseAssets) == 0 {
		baseCurrencyDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
		if found {
			params.BaseAssets = []string{baseCurrencyDenom}
			k.SetParams(ctx, params)
		}
	}

	found := false
	for _, asset := range params.BaseAssets {
		if asset == denom {
			found = true
			break
		}
	}
	return found
}
