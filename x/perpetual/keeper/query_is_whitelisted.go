package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) IsWhitelisted(goCtx context.Context, req *types.IsWhitelistedRequest) (*types.IsWhitelistedResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	accAddress, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}

	isWhitelisted := k.CheckIfWhitelisted(ctx, accAddress)

	return &types.IsWhitelistedResponse{
		Address:       req.Address,
		IsWhitelisted: isWhitelisted,
	}, nil
}
