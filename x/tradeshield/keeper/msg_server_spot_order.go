package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreateSpotOrder(goCtx context.Context, msg *types.MsgCreateSpotOrder) (*types.MsgCreateSpotOrderResponse, error) {
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

	// if the order is market buy, execute it immediately
	if msg.OrderType == types.SpotOrderType_MARKETBUY {
		err := k.ExecuteMarketBuyOrder(ctx, pendingSpotOrder)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgCreateSpotOrderResponse{
		OrderId: id,
	}, nil
}

func (k msgServer) UpdateSpotOrder(goCtx context.Context, msg *types.MsgUpdateSpotOrder) (*types.MsgUpdateSpotOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	// _ := sdk.UnwrapSDKContext(goCtx)

	return &types.MsgUpdateSpotOrderResponse{}, nil
}

func (k msgServer) CancelSpotOrder(goCtx context.Context, msg *types.MsgCancelSpotOrder) (*types.MsgCancelSpotOrderResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the spot order
	spotOrder, found := k.GetPendingSpotOrder(ctx, msg.OrderId)
	if !found {
		return nil, types.ErrSpotOrderNotFound
	}

	if spotOrder.OwnerAddress != msg.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	k.RemovePendingSpotOrder(ctx, msg.OrderId)
	types.EmitCloseSpotOrderEvent(ctx, spotOrder)

	return &types.MsgCancelSpotOrderResponse{
		OrderId: spotOrder.OrderId,
	}, nil
}

func (k msgServer) CancelSpotOrders(goCtx context.Context, msg *types.MsgCancelSpotOrders) (*types.MsgCancelSpotOrdersResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	if len(msg.SpotOrderIds) == 0 {
		return nil, types.ErrSizeZero
	}
	// loop through the spot orders and execute them
	for _, spotOrderId := range msg.SpotOrderIds {
		_, err := k.CancelSpotOrder(goCtx, &types.MsgCancelSpotOrder{
			OwnerAddress: msg.Creator,
			OrderId:      spotOrderId,
		})
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgCancelSpotOrdersResponse{}, nil
}
