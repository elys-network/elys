package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) MTP(goCtx context.Context, req *types.MTPRequest) (*types.MTPResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	mtp, err := k.GetMTP(ctx, req.Address, req.Id)
	if err != nil {
		return &types.MTPResponse{}, nil
	}

	return &types.MTPResponse{Mtp: &mtp}, nil
}
