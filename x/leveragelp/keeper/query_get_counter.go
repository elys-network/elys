package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) QueryPositionCounter(goCtx context.Context, req *types.PositionCounterRequest) (*types.PositionCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.PositionCounterResponse{
		Result: k.GetPositionCounter(ctx, req.PoolId),
	}, nil
}
