package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/stablestake/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AmmPool(goCtx context.Context, req *types.QueryAmmPoolRequest) (*types.QueryAmmPoolResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAmmPoolResponse{AmmPool: k.GetAmmPool(ctx, req.Id)}, nil
}

func (k Keeper) AllAmmPools(goCtx context.Context, req *types.QueryAllAmmPoolsRequest) (*types.QueryAllAmmPoolsResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	return &types.QueryAllAmmPoolsResponse{AmmPools: k.GetAllAmmPools(ctx)}, nil
}
