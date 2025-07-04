package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetCounter(goCtx context.Context, req *types.CounterRequest) (*types.CounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.CounterResponse{
		Result: k.GetPositionCounter(ctx, req.PoolId),
	}, nil
}
