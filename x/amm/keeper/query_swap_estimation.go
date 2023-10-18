package keeper

import (
	"context"

    "github.com/elys-network/elys/x/amm/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) SwapEstimation(goCtx context.Context,  req *types.QuerySwapEstimationRequest) (*types.QuerySwapEstimationResponse, error) {
	if req == nil {
        return nil, status.Error(codes.InvalidArgument, "invalid request")
    }

	ctx := sdk.UnwrapSDKContext(goCtx)

    // TODO: Process the query
    _ = ctx

	return &types.QuerySwapEstimationResponse{}, nil
}
