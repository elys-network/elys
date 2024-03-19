package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Bonus(goCtx context.Context, req *types.QueryBonusRequest) (*types.QueryBonusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryBonusResponse{}, nil
}

func (k Keeper) BuyElysEst(goCtx context.Context, req *types.QueryBuyElysEstRequest) (*types.QueryBuyElysEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryBuyElysEstResponse{}, nil
}

func (k Keeper) ReturnElysEst(goCtx context.Context, req *types.QueryReturnElysEstRequest) (*types.QueryReturnElysEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryReturnElysEstResponse{}, nil
}

func (k Keeper) Orders(goCtx context.Context, req *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryOrdersResponse{}, nil
}

func (k Keeper) AllOrders(goCtx context.Context, req *types.QueryAllOrdersRequest) (*types.QueryAllOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	_ = ctx
	// TODO: implement

	return &types.QueryAllOrdersResponse{}, nil
}
