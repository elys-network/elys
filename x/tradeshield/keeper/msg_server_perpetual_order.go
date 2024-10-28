package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePerpetualOpenOrder(goCtx context.Context, msg *types.MsgCreatePerpetualOpenOrder) (*types.MsgCreatePerpetualOpenOrderResponse, error) {
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

	return &types.MsgCreatePerpetualOpenOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) CreatePerpetualCloseOrder(goCtx context.Context, msg *types.MsgCreatePerpetualCloseOrder) (*types.MsgCreatePerpetualCloseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingPerpetualOrder = types.PerpetualOrder{
		PerpetualOrderType: msg.OrderType,
		TriggerPrice:       msg.TriggerPrice,
		OwnerAddress:       msg.OwnerAddress,
		PositionId:         msg.PositionId,
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	return &types.MsgCreatePerpetualCloseOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) UpdatePerpetualOrder(goCtx context.Context, msg *types.MsgUpdatePerpetualOrder) (*types.MsgUpdatePerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingPerpetualOrder = types.PerpetualOrder{
		OwnerAddress: msg.OwnerAddress,
	}

	// Checks that the element exists
	val, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != val.OwnerAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.SetPendingPerpetualOrder(ctx, pendingPerpetualOrder)

	return &types.MsgUpdatePerpetualOrderResponse{}, nil
}

func (k msgServer) CancelPerpetualOrders(goCtx context.Context, msg *types.MsgCancelPerpetualOrders) (*types.MsgCancelPerpetualOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	if len(msg.OrderIds) == 0 {
		return nil, types.ErrSizeZero
	}
	// loop through the spot orders and execute them
	for _, OrderId := range msg.OrderIds {

		// Checks that the element exists
		val, found := k.GetPendingPerpetualOrder(ctx, OrderId)
		if !found {
			return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", OrderId))
		}

		// Checks if the msg creator is the same as the current owner
		if msg.OwnerAddress != val.OwnerAddress {
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
		}

		k.RemovePendingPerpetualOrder(ctx, OrderId)
		types.EmitClosePerpetualOrderEvent(ctx, val)
	}
	return &types.MsgCancelPerpetualOrdersResponse{}, nil
}
