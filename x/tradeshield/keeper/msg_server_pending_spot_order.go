package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePendingSpotOrder(goCtx context.Context, msg *types.MsgCreatePendingSpotOrder) (*types.MsgCreatePendingSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingSpotOrder = types.SpotOrder{
		OrderType:    msg.OrderType,
		OrderId:      uint64(0),
		OrderPrice:   msg.OrderPrice,
		OrderAmount:  *msg.OrderAmount,
		OwnerAddress: msg.OwnerAddress,
	}

	id := k.AppendPendingSpotOrder(
		ctx,
		pendingSpotOrder,
	)

	return &types.MsgCreatePendingSpotOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) UpdatePendingSpotOrder(goCtx context.Context, msg *types.MsgUpdatePendingSpotOrder) (*types.MsgUpdatePendingSpotOrderResponse, error) {
	// _ := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdatePendingSpotOrderResponse{}, nil
}

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
		types.EmitCloseSpotOrderEvent(ctx, spotOrder)
	}

	return &types.MsgCancelSpotOrdersResponse{}, nil
}