package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/masterchef/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) Apr(goCtx context.Context, req *types.QueryAprRequest) (*types.QueryAprResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	apr, err := k.CalculateApr(ctx, req)
	if err != nil {
		return nil, err
	}

	return &types.QueryAprResponse{Apr: apr.Dec()}, nil
}
