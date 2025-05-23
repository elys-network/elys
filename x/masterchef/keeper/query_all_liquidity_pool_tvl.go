package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/masterchef/types"
	"github.com/osmosis-labs/osmosis/osmomath"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllLiquidityPoolTVL(goCtx context.Context, req *types.QueryAllLiquidityPoolTVLRequest) (*types.QueryAllLiquidityPoolTVLResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	allPools := k.amm.GetAllPool(ctx)
	poolsTVL := osmomath.ZeroBigDec()
	totalTVL := math.ZeroInt()

	for _, pool := range allPools {
		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}

		poolsTVL = poolsTVL.Add(tvl)
	}
	totalTVL = totalTVL.Add(poolsTVL.Dec().TruncateInt())

	stableStakeTVL := k.stableKeeper.AllTVL(ctx)
	totalTVL = totalTVL.Add(stableStakeTVL.Dec().TruncateInt())

	return &types.QueryAllLiquidityPoolTVLResponse{
		Total:       totalTVL,
		Pools:       poolsTVL.Dec().TruncateInt(),
		UsdcStaking: stableStakeTVL.Dec().TruncateInt(),
	}, nil
}
