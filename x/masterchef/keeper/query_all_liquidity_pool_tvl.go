package keeper

import (
	"context"

	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	elystypes "github.com/elys-network/elys/types"
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
	poolsTVL := elystypes.ZeroDec34()
	totalTVL := math.ZeroInt()

	for _, pool := range allPools {
		tvl, err := pool.TVL(ctx, k.oracleKeeper, k.accountedPoolKeeper)
		if err != nil {
			return nil, err
		}

		poolsTVL = poolsTVL.Add(tvl)
	}
	totalTVL = totalTVL.Add(poolsTVL.ToInt())

	stableStakeTVL := k.stableKeeper.AllTVL(ctx, k.oracleKeeper)
	totalTVL = totalTVL.Add(stableStakeTVL.ToInt())

	return &types.QueryAllLiquidityPoolTVLResponse{
		Total:       totalTVL,
		Pools:       poolsTVL.ToInt(),
		UsdcStaking: stableStakeTVL.ToInt(),
	}, nil
}
