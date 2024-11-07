package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	perpetualtypes "github.com/elys-network/elys/x/perpetual/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreatePerpetualOpenOrder(goCtx context.Context, msg *types.MsgCreatePerpetualOpenOrder) (*types.MsgCreatePerpetualOpenOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Verify if perpetual pool exists
	_, found := k.perpetual.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("pool %d not found", msg.PoolId))
	}

	var pendingPerpetualOrder = types.PerpetualOrder{
		PerpetualOrderType: types.PerpetualOrderType_LIMITOPEN,
		TriggerPrice:       msg.TriggerPrice,
		Collateral:         msg.Collateral,
		OwnerAddress:       msg.OwnerAddress,
		TradingAsset:       msg.TradingAsset,
		Position:           msg.Position,
		Leverage:           msg.Leverage,
		TakeProfitPrice:    msg.TakeProfitPrice,
		StopLossPrice:      msg.StopLossPrice,
		PoolId:             msg.PoolId,
		PositionId:         0,
	}

	id := k.AppendPendingPerpetualOrder(
		ctx,
		pendingPerpetualOrder,
	)

	// Verify if order is valid before saving
	_, err := k.perpetual.HandleOpenEstimation(ctx, &perpetualtypes.QueryOpenEstimationRequest{
		Position:        perpetualtypes.Position(msg.Position),
		Leverage:        msg.Leverage,
		TradingAsset:    msg.TradingAsset,
		Collateral:      msg.Collateral,
		TakeProfitPrice: msg.TakeProfitPrice,
		PoolId:          msg.PoolId,
	})
	if err != nil {
		return nil, err
	}

	return &types.MsgCreatePerpetualOpenOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) CreatePerpetualCloseOrder(goCtx context.Context, msg *types.MsgCreatePerpetualCloseOrder) (*types.MsgCreatePerpetualCloseOrderResponse, error) {
	// Disable for v1
	return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "disabled for v1")
	// ctx := sdk.UnwrapSDKContext(goCtx)

	// // check if the position owner address matches the msg owner address
	// position, err := k.perpetual.GetMTP(ctx, sdk.MustAccAddressFromBech32(msg.OwnerAddress), msg.PositionId)
	// if err != nil {
	// 	return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("position %d not found", msg.PositionId))
	// }

	// var pendingPerpetualOrder = types.PerpetualOrder{
	// 	PerpetualOrderType: types.PerpetualOrderType_LIMITCLOSE,
	// 	TriggerPrice: &types.TriggerPrice{
	// 		TradingAssetDenom: position.TradingAsset,
	// 		Rate:              msg.TriggerPrice.Rate,
	// 	},
	// 	OwnerAddress: position.Address,
	// 	PositionId:   position.Id,
	// }

	// id := k.AppendPendingPerpetualOrder(
	// 	ctx,
	// 	pendingPerpetualOrder,
	// )

	// return &types.MsgCreatePerpetualCloseOrderResponse{
	// 	OrderId: id,
	// }, nil
}

func (k msgServer) UpdatePerpetualOrder(goCtx context.Context, msg *types.MsgUpdatePerpetualOrder) (*types.MsgUpdatePerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	order, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("key %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != order.OwnerAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	order.TriggerPrice = msg.TriggerPrice
	k.SetPendingPerpetualOrder(ctx, order)

	return &types.MsgUpdatePerpetualOrderResponse{}, nil
}

func (k msgServer) CancelPerpetualOrder(goCtx context.Context, msg *types.MsgCancelPerpetualOrder) (*types.MsgCancelPerpetualOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	order, found := k.GetPendingPerpetualOrder(ctx, msg.OrderId)
	if !found {
		return nil, sdkerrors.Wrap(sdkerrors.ErrKeyNotFound, fmt.Sprintf("order %d doesn't exist", msg.OrderId))
	}

	// Checks if the msg creator is the same as the current owner
	if msg.OwnerAddress != order.OwnerAddress {
		return nil, sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingPerpetualOrder(ctx, msg.OrderId)
	types.EmitCancelPerpetualOrderEvent(ctx, order)

	return &types.MsgCancelPerpetualOrderResponse{
		OrderId: order.OrderId,
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
