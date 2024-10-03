package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

// TODO: Complete message in another task
func (k msgServer) CreatePendingPerpetualOrder(goCtx context.Context, msg *types.MsgCreatePendingPerpetualOrder) (*types.MsgCreatePendingPerpetualOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingPerpetualOrder = types.PerpetualOrder{
		PerpetualOrderType: msg.OrderType,
		TriggerPrice:       msg.TriggerPrice,
		Collateral:         msg.Collateral,
		OwnerAddress:       msg.OwnerAddress,
		TradingAsset:       msg.TradingAsset,
		Position:           msg.Position,
		Leverage:           msg.Leverage,
		TakeProfitPrice:    msg.TakeProfitPrice,
		StopLossPrice:      msg.StopLossPrice,
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	return &types.MsgCreatePendingPerpetualOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) UpdatePendingPerpetualOrder(goCtx context.Context, msg *types.MsgUpdatePendingPerpetualOrder) (*types.MsgUpdatePendingPerpetualOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingPerpetualOrder = types.PerpetualOrder{
		OwnerAddress: msg.OwnerAddress,
	}

	// Checks that the element exists
	val, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != val.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPendingPerpetualOrder(ctx, pendingPerpetualOrder)

	return &types.MsgUpdatePendingPerpetualOrderResponse{}, nil
}

func (k msgServer) DeletePendingPerpetualOrder(goCtx context.Context, msg *types.MsgDeletePendingPerpetualOrder) (*types.MsgDeletePendingPerpetualOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, errorsmod.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != val.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingPerpetualOrder(ctx, msg.OrderId)

	return &types.MsgDeletePendingPerpetualOrderResponse{}, nil
}
