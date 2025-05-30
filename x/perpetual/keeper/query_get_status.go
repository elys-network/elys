package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) GetStatus(goCtx context.Context, req *types.StatusRequest) (*types.StatusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	return &types.StatusResponse{
		OpenMtpCount:     k.GetOpenMTPCount(ctx),
		LifetimeMtpCount: k.GetMTPCount(ctx),
	}, nil
}
