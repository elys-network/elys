package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) IsWhitelisted(goCtx context.Context, req *types.IsWhitelistedRequest) (*types.IsWhitelistedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	isWhitelisted := k.CheckIfWhitelisted(ctx, req.Address)

	return &types.IsWhitelistedResponse{
		Address:       req.Address,
		IsWhitelisted: isWhitelisted,
	}, nil
}
