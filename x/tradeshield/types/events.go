package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const (
	TypeEvtCloseSpotOrder      = "close_spot_order"
	TypeEvtClosePerpetualOrder = "close_perpetual_order"
)

func EmitCloseSpotOrderEvent(ctx sdk.Context, order SpotOrder) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewCloseSpotOrderEvt(order),
	})
}

func NewCloseSpotOrderEvt(order SpotOrder) sdk.Event {
	return sdk.NewEvent(TypeEvtCloseSpotOrder,
		sdk.NewAttribute("order_type", order.OrderType.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("order_target_denom", order.OrderTargetDenom),
		sdk.NewAttribute("status", order.Status.String()),
		sdk.NewAttribute("order_price", order.OrderPrice.String()),
		sdk.NewAttribute("order_amount", order.OrderAmount.String()),
		sdk.NewAttribute("date", order.Date.String()),
	)
}

func EmitClosePerpetualOrderEvent(ctx sdk.Context, order PerpetualOrder) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewClosePerpetualOrderEvt(order),
	})
}

func NewClosePerpetualOrderEvt(order PerpetualOrder) sdk.Event {
	return sdk.NewEvent(TypeEvtCloseSpotOrder,
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("position", order.Position.String()),
		sdk.NewAttribute("trigger_price", order.TriggerPrice.String()),
		sdk.NewAttribute("collateral", order.Collateral.String()),
		sdk.NewAttribute("trading_asset", order.TradingAsset),
		sdk.NewAttribute("leverage", order.Leverage.String()),
		sdk.NewAttribute("take_profit_price", order.TakeProfitPrice.String()),
		sdk.NewAttribute("position_id", strconv.FormatInt(int64(order.PositionId), 10)),
		sdk.NewAttribute("status", order.Status.String()),
		sdk.NewAttribute("stop_loss_price", order.StopLossPrice.String()),
	)
}
