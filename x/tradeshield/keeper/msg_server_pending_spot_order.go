package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePendingSpotOrder(goCtx context.Context, msg *types.MsgCreatePendingSpotOrder) (*types.MsgCreatePendingSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingSpotOrder = types.PendingSpotOrder{
		Creator: msg.Creator,
		Order:   msg.Order,
	}

	id := k.AppendPendingSpotOrder(
		ctx,
		pendingSpotOrder,
	)

	return &types.MsgCreatePendingSpotOrderResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdatePendingSpotOrder(goCtx context.Context, msg *types.MsgUpdatePendingSpotOrder) (*types.MsgUpdatePendingSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingSpotOrder = types.PendingSpotOrder{
		Creator: msg.Creator,
		Id:      msg.Id,
		Order:   msg.Order,
	}

	// Checks that the element exists
	val, found := k.GetPendingSpotOrder(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPendingSpotOrder(ctx, pendingSpotOrder)

	return &types.MsgUpdatePendingSpotOrderResponse{}, nil
}

func (k msgServer) DeletePendingSpotOrder(goCtx context.Context, msg *types.MsgDeletePendingSpotOrder) (*types.MsgDeletePendingSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPendingSpotOrder(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.Creator {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingSpotOrder(ctx, msg.Id)

	return &types.MsgDeletePendingSpotOrderResponse{}, nil
}
