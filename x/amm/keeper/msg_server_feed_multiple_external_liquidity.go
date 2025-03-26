package keeper

import (
	"context"

	sdkmath "cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	oracletypes "github.com/ojo-network/ojo/x/oracle/types"
)

func (k Keeper) GetExternalLiquidityRatio(ctx sdk.Context, pool types.Pool, amountDepthInfo []types.AssetAmountDepth) ([]types.PoolAsset, error) {
	updatedAssets := make([]types.PoolAsset, len(pool.PoolAssets))
	copy(updatedAssets, pool.PoolAssets)

	for i, asset := range updatedAssets {
		for _, el := range amountDepthInfo {
			entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, asset.Token.Denom)
			if !found {
				return nil, assetprofiletypes.ErrAssetProfileNotFound
			}
			if entry.DisplayName == el.Asset {

				O_Tvl := el.Amount
				P_Tvl := asset.Token.Amount.ToLegacyDec()

				// Ensure tvl is not zero to avoid division by zero
				if P_Tvl.IsZero() {
					return nil, types.ErrAmountTooLow
				}

				liquidityRatio := LiquidityRatioFromPriceDepth(el.Depth)
				// Ensure tvl is not zero to avoid division by zero
				if liquidityRatio.IsZero() {
					return nil, types.ErrAmountTooLow
				}
				asset.ExternalLiquidityRatio = (O_Tvl.Quo(P_Tvl)).Quo(liquidityRatio)

				if asset.ExternalLiquidityRatio.LT(sdkmath.LegacyOneDec()) {
					asset.ExternalLiquidityRatio = sdkmath.LegacyOneDec()
				}
			}
		}
		updatedAssets[i] = asset
	}
	return updatedAssets, nil
}

func LiquidityRatioFromPriceDepth(depth sdkmath.LegacyDec) sdkmath.LegacyDec {
	if depth == sdkmath.LegacyOneDec() {
		return sdkmath.LegacyOneDec()
	}
	sqrt, err := sdkmath.LegacyOneDec().Sub(depth).ApproxSqrt()
	if err != nil {
		panic(err)
	}
	return sdkmath.LegacyOneDec().Sub(sqrt)
}

func (k msgServer) FeedMultipleExternalLiquidity(goCtx context.Context, msg *types.MsgFeedMultipleExternalLiquidity) (*types.MsgFeedMultipleExternalLiquidityResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	feeder, found := k.oracleKeeper.GetPriceFeeder(ctx, sdk.MustAccAddressFromBech32(msg.Sender))
	if !found {
		return nil, oracletypes.ErrNotAPriceFeeder
	}

	if !feeder.IsActive {
		return nil, oracletypes.ErrPriceFeederNotActive
	}

	for _, el := range msg.Liquidity {
		pool, found := k.GetPool(ctx, el.PoolId)
		if !found {
			return nil, types.ErrInvalidPoolId
		}

		// Get external liquidity ratio for each of the asset separately
		poolAssets, err := k.GetExternalLiquidityRatio(ctx, pool, el.AmountDepthInfo)
		if err != nil {
			return nil, err
		}

		pool.PoolAssets = poolAssets
		k.SetPool(ctx, pool)
	}

	return &types.MsgFeedMultipleExternalLiquidityResponse{}, nil
}
