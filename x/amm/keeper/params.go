package keeper

import (
	"slices"

	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v7/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v7/x/parameter/types"
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

	found := slices.Contains(params.BaseAssets, denom)
	return found
}

func (k Keeper) V8Migrate(ctx sdk.Context) error {
	baseCurrencyDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	params := types.DefaultParams()
	params.BaseAssets = []string{baseCurrencyDenom}

	k.SetParams(ctx, params)

	legacyPools := k.GetAllLegacyPool(ctx)
	for _, legacyPool := range legacyPools {
		var newPool types.Pool
		newPool.PoolId = legacyPool.PoolId
		newPool.Address = legacyPool.Address
		newPool.PoolParams = types.PoolParams{
			SwapFee:   legacyPool.PoolParams.SwapFee,
			UseOracle: legacyPool.PoolParams.UseOracle,
			FeeDenom:  legacyPool.PoolParams.FeeDenom,
		}
		newPool.TotalShares = legacyPool.TotalShares
		newPool.TotalWeight = legacyPool.TotalWeight
		newPool.PoolAssets = legacyPool.PoolAssets
		newPool.RebalanceTreasury = legacyPool.RebalanceTreasury

		k.SetPool(ctx, newPool)
	}

	return nil
}
