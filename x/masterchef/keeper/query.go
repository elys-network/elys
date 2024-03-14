package keeper

import (
	"context"

	"github.com/elys-network/elys/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) ExternalIncentive(goCtx context.Context, req *types.QueryExternalIncentiveRequest) (*types.QueryExternalIncentiveResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid request")
}

func (k Keeper) PoolInfo(goCtx context.Context, req *types.QueryPoolInfoRequest) (*types.QueryPoolInfoResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid request")
}

func (k Keeper) PoolRewardInfo(goCtx context.Context, req *types.QueryPoolRewardInfoRequest) (*types.QueryPoolRewardInfoResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid request")
}

func (k Keeper) UserRewardInfo(goCtx context.Context, req *types.QueryUserRewardInfoRequest) (*types.QueryUserRewardInfoResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid request")
}

func (k Keeper) UserPendingReward(goCtx context.Context, req *types.QueryUserPendingRewardRequest) (*types.QueryUserPendingRewardResponse, error) {
	return nil, status.Error(codes.InvalidArgument, "invalid request")
}
