package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) LockedOrder(goCtx context.Context, req *types.QueryLockedOrderRequest) (*types.QueryLockedOrderResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(req.User)
	total := k.RetrieveTradeshieldTotal(ctx, sender)

	return &types.QueryLockedOrderResponse{
		Total: total.String(),
	}, nil
}
