package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

// UpdatePoolParams updates the pool params
func (k Keeper) UpdatePoolParams(ctx sdk.Context, poolId uint64, poolParams types.PoolParams) (uint64, types.PoolParams, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return 0, types.PoolParams{}, types.ErrPoolNotFound
	}

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return 0, types.PoolParams{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}
	baseCurrency := entry.Denom

	// If the fee denom is empty, set it to the base currency
	if poolParams.FeeDenom == "" {
		poolParams.FeeDenom = baseCurrency
	}

	pool.PoolParams = poolParams
	k.SetPool(ctx, pool)

	return pool.PoolId, pool.PoolParams, nil
}
