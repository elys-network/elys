package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) CreateSpotOrder(goCtx context.Context, msg *types.MsgCreateSpotOrder) (*types.MsgCreateSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	var pendingSpotOrder = types.SpotOrder{
		OrderType:        msg.OrderType,
		OrderId:          uint64(0),
		OrderPrice:       msg.OrderPrice,
		OrderAmount:      msg.OrderAmount,
		OwnerAddress:     msg.OwnerAddress,
		OrderTargetDenom: msg.OrderTargetDenom,
		Date:             &types.Date{Height: uint64(ctx.BlockHeight()), Timestamp: uint64(ctx.BlockTime().Unix())},
	}

	// if the order is market buy, execute it immediately
	if msg.OrderType == types.SpotOrderType_MARKETBUY {
		res, err := k.ExecuteMarketBuyOrder(ctx, pendingSpotOrder)
		if err != nil {
			return nil, err
		}

		ctx.EventManager().EmitEvent(types.NewExecuteMarketBuySpotOrderEvt(pendingSpotOrder, res))

		return &types.MsgCreateSpotOrderResponse{
			OrderId: pendingSpotOrder.OrderId,
		}, nil
	}

	// add the order to the pending orders
	id := k.AppendPendingSpotOrder(
		ctx,
		pendingSpotOrder,
	)

	// set the order id
	pendingSpotOrder.OrderId = id

	// send order amount from owner to the order address
	ownerAddress := sdk.MustAccAddressFromBech32(pendingSpotOrder.OwnerAddress)
	err := k.Keeper.bank.SendCoins(ctx, ownerAddress, pendingSpotOrder.GetOrderAddress(), sdk.NewCoins(pendingSpotOrder.OrderAmount))
	if err != nil {
		return nil, err
	}

	// return the order id
	return &types.MsgCreateSpotOrderResponse{
		OrderId: pendingSpotOrder.OrderId,
	}, nil
}

func (k msgServer) UpdateSpotOrder(goCtx context.Context, msg *types.MsgUpdateSpotOrder) (*types.MsgUpdateSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Checks that the element exists
	order, found := k.GetPendingSpotOrder(ctx, msg.OrderId)
	if !found {
		return nil, types.ErrSpotOrderNotFound
	}

	// Checks if the the msg creator is the same as the current owner
	if msg.OwnerAddress != order.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// update the order
	order.OrderPrice = msg.OrderPrice
	k.SetPendingSpotOrder(ctx, order)

	return &types.MsgUpdateSpotOrderResponse{}, nil
}

func (k msgServer) CancelSpotOrder(goCtx context.Context, msg *types.MsgCancelSpotOrder) (*types.MsgCancelSpotOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// get the spot order
	spotOrder, found := k.GetPendingSpotOrder(ctx, msg.OrderId)
	if !found {
		return nil, types.ErrSpotOrderNotFound
	}

	if spotOrder.OwnerAddress != msg.OwnerAddress {
		return nil, errorsmod.Wrap(sdkerrors.ErrUnauthorized, "incorrect owner")
	}

	// Get all balances from the spot order address
	orderAddress := spotOrder.GetOrderAddress()
	balances := k.Keeper.bank.GetAllBalances(ctx, orderAddress)

	// Send all available balances back to the owner if there are any
	if !balances.IsZero() {
		ownerAddress := sdk.MustAccAddressFromBech32(spotOrder.OwnerAddress)
		err := k.Keeper.bank.SendCoins(ctx, orderAddress, ownerAddress, balances)
		if err != nil {
			return nil, err
		}
	}

	k.RemovePendingSpotOrder(ctx, msg.OrderId)
	types.EmitCloseSpotOrderEvent(ctx, spotOrder)

	return &types.MsgCancelSpotOrderResponse{
		OrderId: spotOrder.OrderId,
	}, nil
}

func (k msgServer) CancelSpotOrders(goCtx context.Context, msg *types.MsgCancelSpotOrders) (*types.MsgCancelSpotOrdersResponse, error) {
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
