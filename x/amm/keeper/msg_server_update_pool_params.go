package keeper

import (
	"context"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	govtypes "github.com/cosmos/cosmos-sdk/x/gov/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/v5/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/v5/x/parameter/types"
)

// UpdatePoolParams updates the pool params
func (k Keeper) UpdatePoolParams(ctx sdk.Context, poolId uint64, newPoolParams types.PoolParams) (uint64, types.PoolParams, error) {
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return 0, types.PoolParams{}, types.ErrPoolNotFound
	}

	baseCurrency, found := k.assetProfileKeeper.GetUsdcDenom(ctx)
	if !found {
		return 0, types.PoolParams{}, errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency)
	}

	// If the fee denom is empty, set it to the base currency
	if newPoolParams.FeeDenom == "" {
		newPoolParams.FeeDenom = baseCurrency
	}

	// changing from non-oracle pool to oracle pool
	if !pool.PoolParams.UseOracle && newPoolParams.UseOracle {

		nonBaseCurrencyDenom := ""
		usdcDenomFound := false

		for _, asset := range pool.PoolAssets {
			if asset.Token.Denom != baseCurrency {
				nonBaseCurrencyDenom = asset.Token.Denom
			}
			if asset.Token.Denom == baseCurrency {
				usdcDenomFound = true
			}
		}
		if !usdcDenomFound {
			return 0, types.PoolParams{}, fmt.Errorf("no usdc denom in the amm pool %d", poolId)
		}
		if nonBaseCurrencyDenom == "" {
			return 0, types.PoolParams{}, fmt.Errorf("no non-usdc denom in the amm pool %d", poolId)
		}

		entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, nonBaseCurrencyDenom)
		if !found {
			return 0, types.PoolParams{}, fmt.Errorf("asset profile for %s not found", nonBaseCurrencyDenom)
		}
		_, found = k.oracleKeeper.GetAssetPrice(ctx, entry.DisplayName)
		if !found {
			return 0, types.PoolParams{}, fmt.Errorf("oracle price for %s not found", entry.DisplayName)
		}
	}
	pool.PoolParams = newPoolParams
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
