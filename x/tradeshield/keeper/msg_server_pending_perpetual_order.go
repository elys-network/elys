package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// TODO: Complete message in another task
func (k msgServer) CreatePendingPerpetualOrder(goCtx context.Context, msg *types.MsgCreatePendingPerpetualOrder) (*types.MsgCreatePendingPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingPerpetualOrder = types.PerpetualOrder{
		OwnerAddress: msg.Creator,
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	return &types.MsgCreatePendingPerpetualOrderResponse{
		Id: id,
	}, nil
}

func (k msgServer) UpdatePendingPerpetualOrder(goCtx context.Context, msg *types.MsgUpdatePendingPerpetualOrder) (*types.MsgUpdatePendingPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingPerpetualOrder = types.PerpetualOrder{
		OwnerAddress: msg.Creator,
	}

	// Checks that the element exists
	val, found := k.GetPendingPerpetualOrder(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.OwnerAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPendingPerpetualOrder(ctx, pendingPerpetualOrder)

	return &types.MsgUpdatePendingPerpetualOrderResponse{}, nil
}

func (k msgServer) DeletePendingPerpetualOrder(goCtx context.Context, msg *types.MsgDeletePendingPerpetualOrder) (*types.MsgDeletePendingPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPendingPerpetualOrder(ctx, msg.Id)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.Id))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.Creator != val.OwnerAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingPerpetualOrder(ctx, msg.Id)

	return &types.MsgDeletePendingPerpetualOrderResponse{}, nil
}
