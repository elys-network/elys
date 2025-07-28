package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetStatus(goCtx context.Context, req *types.StatusRequest) (*types.StatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.StatusResponse{
		OpenPositionCount:     k.GetOpenPositionCount(ctx),
		LifetimePositionCount: k.GetPositionCount(ctx),
	}, nil
}
