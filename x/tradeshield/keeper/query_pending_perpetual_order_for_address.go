package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/tradeshield/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PendingPerpetualOrderForAddress(goCtx context.Context, req *types.QueryPendingPerpetualOrderForAddressRequest) (*types.QueryPendingPerpetualOrderForAddressResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)

	pendingPerpetualOrders, _, err := k.GetPendingPerpetualOrdersForAddress(ctx, req.Address, &req.Status, nil)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pendingPerpetualOrdersExtraInfo []types.PerpetualOrderExtraInfo

	for _, perpetualOrder := range pendingPerpetualOrders {
		pendingPerpetualOrderExtraInfo, err := k.ConstructPerpetualOrderExtraInfo(ctx, perpetualOrder)
		if err != nil {
			return nil, err
		}

		pendingPerpetualOrdersExtraInfo = append(pendingPerpetualOrdersExtraInfo, *pendingPerpetualOrderExtraInfo)
	}

	return &types.QueryPendingPerpetualOrderForAddressResponse{
		PendingPerpetualOrders: pendingPerpetualOrdersExtraInfo,
	}, nil
}
