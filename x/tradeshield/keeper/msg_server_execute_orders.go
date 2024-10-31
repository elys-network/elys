package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) ExecuteOrders(goCtx context.Context, msg *types.MsgExecuteOrders) (*types.MsgExecuteOrdersResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// loop through the spot orders and execute them
	for _, spotOrderId := range msg.SpotOrderIds {
		// get the spot order
		spotOrder, found := k.GetPendingSpotOrder(ctx, spotOrderId)
		if !found {
			return nil, types.ErrSpotOrderNotFound
		}

		var err error

		// dispatch based on the order type
		switch spotOrder.OrderType {
		case types.SpotOrderType_STOPLOSS:
			// execute the stop loss order
			err = k.ExecuteStopLossOrder(ctx, spotOrder)
		case types.SpotOrderType_LIMITSELL:
			// execute the limit sell order
			err = k.ExecuteLimitSellOrder(ctx, spotOrder)
		case types.SpotOrderType_LIMITBUY:
			// execute the limit buy order
			err = k.ExecuteLimitBuyOrder(ctx, spotOrder)
		case types.SpotOrderType_MARKETBUY:
			// execute the market buy order
			err = k.ExecuteMarketBuyOrder(ctx, spotOrder)
		}

		// return the error if any
		if err != nil {
			return nil, err
		}
	}

	// loop through the perpetual orders and execute them
	for _, perpetualOrderId := range msg.PerpetualOrderIds {
		// get the perpetual order
		perpetualOrder, found := k.GetPendingPerpetualOrder(ctx, perpetualOrderId)
		if !found {
			return nil, types.ErrPerpetualOrderNotFound
		}

		var err error

		// dispatch based on the order type
		switch perpetualOrder.PerpetualOrderType {
		case types.PerpetualOrderType_LIMITOPEN:
			// execute the limit open order
			err = k.ExecuteLimitOpenOrder(ctx, perpetualOrder)
		case types.PerpetualOrderType_LIMITCLOSE:
			// execute the limit close order
			err = k.ExecuteLimitCloseOrder(ctx, perpetualOrder)
		}

		// return the error if any
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgExecuteOrdersResponse{}, nil
}
