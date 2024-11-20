package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/tradeshield/types"
)

func (k msgServer) ExecuteOrders(goCtx context.Context, msg *types.MsgExecuteOrders) (*types.MsgExecuteOrdersResponse, error) {
	if err := msg.ValidateBasic(); err != nil {
		return nil, err
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	spotLog := []string{}
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

		// log the error if any
		if err != nil {
			// Add log about error or not executed
			spotLog = append(spotLog, fmt.Sprintf("Spot order Id:%d cannot be executed due to err: %s", spotOrderId, err.Error()))
		}
	}

	perpLog := []string{}
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
			// Disable for v1
			// case types.PerpetualOrderType_LIMITCLOSE:
			// 	// execute the limit close order
			// 	err = k.ExecuteLimitCloseOrder(ctx, perpetualOrder)
		}

		// return the error if any
		// log the error if any
		if err != nil {
			// Add log about error or not executed
			perpLog = append(perpLog, fmt.Sprintf("Perpetual order Id:%d cannot be executed due to err: %s", perpetualOrderId, err.Error()))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.TypeEvtExecuteOrders,
		sdk.NewAttribute("liquidations", strings.Join(spotLog, "\n")),
		sdk.NewAttribute("stop_loss", strings.Join(perpLog, "\n")),
	))

	return &types.MsgExecuteOrdersResponse{}, nil
}
