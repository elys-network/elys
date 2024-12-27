package keeper

import (
	"context"
	"cosmossdk.io/math"
	ptypes "github.com/elys-network/elys/x/parameter/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllLiquidityPoolTVL(goCtx context.Context, req *types.QueryAllLiquidityPoolTVLRequest) (*types.QueryAllLiquidityPoolTVLResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	allPools := k.amm.GetAllPool(ctx)
	poolsTVL := math.LegacyZeroDec()
	totalTVL := math.ZeroInt()

	for _, pool := range allPools {
		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}

		poolsTVL = poolsTVL.Add(tvl)
	}
	totalTVL = totalTVL.Add(poolsTVL.TruncateInt())

	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		return nil, status.Error(codes.NotFound, "asset profile not found")
	}

	stableStakeTVL := k.stableKeeper.TVL(ctx, k.oracleKeeper, entry.Denom)
	totalTVL = totalTVL.Add(stableStakeTVL.TruncateInt())

	return &types.QueryAllLiquidityPoolTVLResponse{
		Total:       totalTVL,
		Pools:       poolsTVL.TruncateInt(),
		UsdcStaking: stableStakeTVL.TruncateInt(),
	}, nil
}
