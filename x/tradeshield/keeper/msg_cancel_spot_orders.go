package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CancelSpotOrders(goCtx context.Context, msg *types.MsgCancelSpotOrders) (*types.MsgCancelSpotOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.SpotOrderIds) == 0 {
		return nil, types.ErrSizeZero
	}
	// loop through the spot orders and execute them
	for _, spotOrderId := range msg.SpotOrderIds {
		// get the spot order
		spotOrder, found := k.GetPendingSpotOrder(ctx, spotOrderId)
		if !found {
			return nil, types.ErrSpotOrderNotFound
		}

		if spotOrder.OwnerAddress != msg.Creator {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
		}
		
		k.RemovePendingSpotOrder(ctx, spotOrderId)
		types.EmitCloseOrdersEvent(ctx, spotOrder)
	}

	return &types.MsgCancelSpotOrdersResponse{}, nil
}