package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/amm/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SlippageTrack(goCtx context.Context, req *types.QuerySlippageTrackRequest) (*types.QuerySlippageTrackResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QuerySlippageTrackResponse{
		Track: types.OraclePoolSlippageTrack{
			PoolId:    req.PoolId,
			Timestamp: uint64(ctx.BlockTime().Unix()),
			Tracked:   k.GetTrackedSlippageDiff(ctx, req.PoolId),
		},
	}, nil
}

func (k Keeper) SlippageTrackAll(goCtx context.Context, req *types.QuerySlippageTrackAllRequest) (*types.QuerySlippageTrackAllResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pools := k.GetAllPool(ctx)
	tracks := []types.OraclePoolSlippageTrack{}
	for _, pool := range pools {
		tracks = append(tracks, types.OraclePoolSlippageTrack{
			PoolId:    pool.PoolId,
			Timestamp: uint64(ctx.BlockTime().Unix()),
			Tracked:   k.GetTrackedSlippageDiff(ctx, pool.PoolId),
		})
	}
	return &types.QuerySlippageTrackAllResponse{
		Tracks: tracks,
	}, nil
}
