package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PerpetualCounter(goCtx context.Context, req *types.PerpetualCounterRequest) (*types.PerpetualCounterResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	v := k.GetPerpetualCounter(ctx, req.Id)
	return &types.PerpetualCounterResponse{
		Result: v,
	}, nil
}
