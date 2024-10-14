package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PendingPerpetualOrderForAddress(goCtx context.Context, req *types.QueryPendingPerpetualOrderForAddressRequest) (*types.QueryPendingPerpetualOrderForAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	res, _, err := k.GetPendingPerpetualOrdersForAddress(ctx, req.Address, &req.Status, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryPendingPerpetualOrderForAddressResponse{
		PendingPerpetualOrders: res,
	}, nil
}
