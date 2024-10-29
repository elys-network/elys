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

	// only allow order type as limit or market open
	if msg.OrderType != types.PerpetualOrderType_LIMITOPEN && msg.OrderType != types.PerpetualOrderType_MARKETOPEN {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid order type")
	}

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
		PoolId:             msg.PoolId,
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	// if the order is market open, execute it immediately
	if msg.OrderType == types.PerpetualOrderType_MARKETOPEN {
		k.ExecuteMarketOpenOrder(ctx, pendingPerpetualOrder)
	}

	return &types.MsgCreatePerpetualOpenOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) CreatePerpetualCloseOrder(goCtx context.Context, msg *types.MsgCreatePerpetualCloseOrder) (*types.MsgCreatePerpetualCloseOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// only allow order type as limit or market close
	if msg.OrderType != types.PerpetualOrderType_LIMITCLOSE && msg.OrderType != types.PerpetualOrderType_MARKETCLOSE {
		return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid order type")
	}

	// check if the position owner address matches the msg owner address
	position, err := k.perpetual.GetMTP(ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), msg.PositionId)
	if err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("position %d not found", msg.PositionId))
	}

	var pendingPerpetualOrder = types.PerpetualOrder{
		PerpetualOrderType: msg.OrderType,
		TriggerPrice: &types.OrderPrice{
			BaseDenom:  position.CollateralAsset,
			QuoteDenom: position.TradingAsset,
			Rate:       msg.TriggerPrice.Rate,
		},
		OwnerAddress: position.Address,
		PositionId:   position.Id,
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	// if the order is market close, execute it immediately
	if msg.OrderType == types.PerpetualOrderType_MARKETCLOSE {
		k.ExecuteMarketCloseOrder(ctx, pendingPerpetualOrder)
	}

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

func (k msgServer) CancelPerpetualOrder(goCtx context.Context, msg *types.MsgCancelPerpetualOrder) (*types.MsgCancelPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	val, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != val.OwnerAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingPerpetualOrder(ctx, msg.OrderId)
	types.EmitClosePerpetualOrderEvent(ctx, val)

	return &types.MsgCancelPerpetualOrderResponse{
		OrderId: val.OrderId,
	}, nil
}

func (k msgServer) CancelPerpetualOrders(goCtx context.Context, msg *types.MsgCancelPerpetualOrders) (*types.MsgCancelPerpetualOrdersResponse, error) {
	if len(msg.OrderIds) == 0 {
		return nil, types.ErrSizeZero
	}
	// loop through the spot orders and cancel them
	for _, orderId := range msg.OrderIds {
		_, err := k.CancelPerpetualOrder(goCtx, &types.MsgCancelPerpetualOrder{
			OwnerAddress: msg.OwnerAddress,
			OrderId:      orderId,
		})
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgCancelPerpetualOrdersResponse{}, nil
}
