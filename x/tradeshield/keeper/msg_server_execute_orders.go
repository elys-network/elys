package keeper

import (
	"context"
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/tradeshield/types"
)

func (k msgServer) ExecuteOrders(goCtx context.Context, msg *types.MsgExecuteOrders) (*types.MsgExecuteOrdersResponse, error) {
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
			cachedCtx, write := ctx.CacheContext()
			_, err = k.ExecuteStopLossOrder(cachedCtx, spotOrder)
			if err == nil {
				write()
			}
		case types.SpotOrderType_LIMITSELL:
			// execute the limit sell order
			cachedCtx, write := ctx.CacheContext()
			_, err = k.ExecuteLimitSellOrder(cachedCtx, spotOrder)
			if err == nil {
				write()
			}
		case types.SpotOrderType_LIMITBUY:
			// execute the limit buy order
			cachedCtx, write := ctx.CacheContext()
			_, err = k.ExecuteLimitBuyOrder(cachedCtx, spotOrder)
			if err == nil {
				write()
			}
		case types.SpotOrderType_MARKETBUY:
			// execute the market buy order
			cachedCtx, write := ctx.CacheContext()
			_, err = k.ExecuteMarketBuyOrder(cachedCtx, spotOrder)
			if err == nil {
				write()
			}
		}

		// log the error if any
		if err != nil {
			// Add log about error or not executed
			spotLog = append(spotLog, fmt.Sprintf("Spot order Id:%d cannot be executed due to err: %s", spotOrderId, err.Error()))
		}
	}

	perpLog := []string{}
	// loop through the perpetual orders and execute them
	for _, perpetualOrderKey := range msg.PerpetualOrders {
		// get the perpetual order
		perpetualOrder, found := k.GetPendingPerpetualOrder(ctx, sdk.MustAccAddressFromBech32(perpetualOrderKey.OwnerAddress), perpetualOrderKey.PoolId, perpetualOrderKey.OrderId)
		if !found {
			return nil, types.ErrPerpetualOrderNotFound
		}

		var err error

		// dispatch based on the order type
		switch perpetualOrder.PerpetualOrderType {
		case types.PerpetualOrderType_LIMITOPEN:
			// execute the limit open order
			cachedCtx, write := ctx.CacheContext()
			err = k.ExecuteLimitOpenOrder(cachedCtx, perpetualOrder)
			if err == nil {
				write()
			}
		case types.PerpetualOrderType_LIMITCLOSE:
			// execute the limit close order
			err = k.ExecuteLimitCloseOrder(ctx, perpetualOrder)
		}

		// return the error if any
		// log the error if any
		if err != nil {
			// Add log about error or not executed
			perpLog = append(perpLog, fmt.Sprintf("Perpetual order Id:%d cannot be executed due to err: %s", perpetualOrder.OrderId, err.Error()))
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.TypeEvtExecuteOrders,
		sdk.NewAttribute("spot_orders", strings.Join(spotLog, "\n")),
		sdk.NewAttribute("perpetual_orders", strings.Join(perpLog, "\n")),
	))

	return &types.MsgExecuteOrdersResponse{}, nil
}
