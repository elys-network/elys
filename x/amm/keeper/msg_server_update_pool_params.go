package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v6/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v6/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v6/x/parameter/types"
)

// UpdatePoolParams updates the pool params
func (k Keeper) UpdatePoolParams(ctx sdk.Context, poolId uint64, newPoolParams types.PoolParams) (uint64, types.PoolParams, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return 0, types.PoolParams{}, types.ErrPoolNotFound
	}

	usdcDenom, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return 0, types.PoolParams{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// If the fee denom is empty, set it to the base currency
	if newPoolParams.FeeDenom == "" {
		newPoolParams.FeeDenom = usdcDenom
	}

	// changing from non-oracle pool to oracle pool
	if !pool.PoolParams.UseOracle && newPoolParams.UseOracle {
		for _, asset := range pool.PoolAssets {
			entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, asset.Token.Denom)
			if !found {
				return 0, types.PoolParams{}, fmt.Errorf("asset profile for %s not found", asset.Token.Denom)
			}
			_, found = k.oracleKeeper.GetAssetPrice(ctx, entry.DisplayName)
			if !found {
				return 0, types.PoolParams{}, fmt.Errorf("oracle price for %s not found", entry.DisplayName)
			}
		}
	}
	pool.PoolParams = newPoolParams
	err := pool.Validate()
	if err != nil {
		return 0, types.PoolParams{}, err
	}
	k.SetPool(ctx, pool)
	return pool.PoolId, pool.PoolParams, nil
}

func (k msgServer) UpdatePoolParams(goCtx context.Context, msg *types.MsgUpdatePoolParams) (*types.MsgUpdatePoolParamsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	poolId, poolParams, err := k.Keeper.UpdatePoolParams(ctx, msg.PoolId, msg.PoolParams)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdatePoolParamsResponse{
		PoolId:     poolId,
		PoolParams: &poolParams,
	}, nil
}
