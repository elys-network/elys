package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PendingSpotOrderForAddress(goCtx context.Context, req *types.QueryPendingSpotOrderForAddressRequest) (*types.QueryPendingSpotOrderForAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	res, _, err := k.GetPendingSpotOrdersForAddress(ctx, req.Address, &req.Status, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPendingSpotOrderForAddressResponse{
		PendingSpotOrders: res,
	}, nil
}
