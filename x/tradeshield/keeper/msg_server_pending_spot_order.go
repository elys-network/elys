package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePendingSpotOrder(goCtx context.Context, msg *types.MsgCreatePendingSpotOrder) (*types.MsgCreatePendingSpotOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
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
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	// _ := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdatePendingSpotOrderResponse{}, nil
}

func (k msgServer) DeletePendingSpotOrder(goCtx context.Context, msg *types.MsgDeletePendingSpotOrder) (*types.MsgDeletePendingSpotOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPendingSpotOrder(ctx, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != val.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingSpotOrder(ctx, msg.OrderId)

	return &types.MsgDeletePendingSpotOrderResponse{}, nil
}
