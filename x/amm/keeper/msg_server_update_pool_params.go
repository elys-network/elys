package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
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

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return 0, types.PoolParams{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// If the fee denom is empty, set it to the base currency
	if poolParams.FeeDenom == "" {
		poolParams.FeeDenom = baseCurrency
	}

	pool.PoolParams = poolParams
	k.SetPool(ctx, pool)
	return pool.PoolId, pool.PoolParams, nil
}

func (k msgServer) UpdatePoolParams(goCtx context.Context, msg *types.MsgUpdatePoolParams) (*types.MsgUpdatePoolParamsResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	if k.authority != msg.Authority {
		return nil, errorsmod.Wrapf(govtypes.ErrInvalidSigner, "invalid authority; expected %s, got %s", k.authority, msg.Authority)
	}

	poolId, poolParams, err := k.Keeper.UpdatePoolParams(ctx, msg.PoolId, *msg.PoolParams)
	if err != nil {
		return nil, err
	}

	return &types.MsgUpdatePoolParamsResponse{
		PoolId:     poolId,
		PoolParams: &poolParams,
	}, nil
}
