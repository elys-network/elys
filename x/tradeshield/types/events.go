package types

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
)

const (
	TypeEvtCloseSpotOrder                  = "tradeshield_close_spot_order"
	TypeEvtCancelPerpetualOrder            = "tradeshield_cancel_perpetual_order"
	TypeEvtExecuteOrders                   = "tradeshield_execute_orders"
	TypeEvtExecuteLimitOpenPerpetualOrder  = "tradeshield_execute_perpetual_limit_open_order"
	TypeEvtExecuteLimitBuySpotOrder        = "tradeshield_execute_limit_buy_spot_order"
	TypeEvtExecuteLimitSellSpotOrder       = "tradeshield_execute_limit_sell_spot_order"
	TypeEvtExecuteStopLossSpotOrder        = "tradeshield_execute_stop_loss_spot_order"
	TypeEvtExecuteMarketBuySpotOrder       = "tradeshield_execute_market_buy_spot_order"
	TypeEvtDeletePendingPerpetualOrder     = "tradeshield_delete_pending_perpetual_order"
	TypeEvtExecuteLimitClosePerpetualOrder = "tradeshield_execute_perpetual_limit_close_order"
	TypeEvtCreatePerpetualOpenOrder        = "tradeshield_create_perpetual_open_order"
	TypeEvtCreatePerpetualCloseOrder       = "tradeshield_create_perpetual_close_order"
	TypeEvtUpdatePerpetualOrder            = "tradeshield_update_perpetual_order"
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

func EmitCancelPerpetualOrderEvent(ctx sdk.Context, order PerpetualOrder) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewCancelPerpetualOrderEvt(order),
	})
}

func NewCancelPerpetualOrderEvt(order PerpetualOrder) sdk.Event {
	return sdk.NewEvent(TypeEvtCancelPerpetualOrder,
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("position", order.Position.String()),
		sdk.NewAttribute("trigger_price", order.TriggerPrice.String()),
		sdk.NewAttribute("collateral", order.Collateral.String()),
		sdk.NewAttribute("leverage", order.Leverage.String()),
		sdk.NewAttribute("take_profit_price", order.TakeProfitPrice.String()),
		sdk.NewAttribute("position_id", strconv.FormatInt(int64(order.PositionId), 10)),
		sdk.NewAttribute("status", order.Status.String()),
		sdk.NewAttribute("stop_loss_price", order.StopLossPrice.String()),
	)
}

func NewExecuteLimitBuySpotOrderEvt(order SpotOrder, res *ammtypes.MsgSwapByDenomResponse) sdk.Event {
	// convert order price to json string
	orderPrice, err := json.Marshal(order.OrderPrice)
	if err != nil {
		panic(err)
	}

	return sdk.NewEvent(TypeEvtExecuteLimitBuySpotOrder,
		sdk.NewAttribute("order_type", order.OrderType.String()),
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("order_price", string(orderPrice)),
		sdk.NewAttribute("order_amount", order.OrderAmount.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_target_denom", order.OrderTargetDenom),
		sdk.NewAttribute("date", order.Date.String()),
		sdk.NewAttribute("amount", res.Amount.String()),
		sdk.NewAttribute("spot_price", res.SpotPrice.String()),
		sdk.NewAttribute("swap_fee", res.SwapFee.String()),
		sdk.NewAttribute("discount", res.Discount.String()),
		sdk.NewAttribute("recipient", res.Recipient),
	)
}

func NewExecuteLimitSellSpotOrderEvt(order SpotOrder, res *ammtypes.MsgSwapByDenomResponse) sdk.Event {
	// convert order price to json string
	orderPrice, err := json.Marshal(order.OrderPrice)
	if err != nil {
		panic(err)
	}

	return sdk.NewEvent(TypeEvtExecuteLimitSellSpotOrder,
		sdk.NewAttribute("order_type", order.OrderType.String()),
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("order_price", string(orderPrice)),
		sdk.NewAttribute("order_amount", order.OrderAmount.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_target_denom", order.OrderTargetDenom),
		sdk.NewAttribute("date", order.Date.String()),
		sdk.NewAttribute("amount", res.Amount.String()),
		sdk.NewAttribute("spot_price", res.SpotPrice.String()),
		sdk.NewAttribute("swap_fee", res.SwapFee.String()),
		sdk.NewAttribute("discount", res.Discount.String()),
		sdk.NewAttribute("recipient", res.Recipient),
	)
}

func NewExecuteStopLossSpotOrderEvt(order SpotOrder, res *ammtypes.MsgSwapByDenomResponse) sdk.Event {
	// convert order price to json string
	orderPrice, err := json.Marshal(order.OrderPrice)
	if err != nil {
		panic(err)
	}

	return sdk.NewEvent(TypeEvtExecuteStopLossSpotOrder,
		sdk.NewAttribute("order_type", order.OrderType.String()),
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("order_price", string(orderPrice)),
		sdk.NewAttribute("order_amount", order.OrderAmount.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_target_denom", order.OrderTargetDenom),
		sdk.NewAttribute("date", order.Date.String()),
		sdk.NewAttribute("amount", res.Amount.String()),
		sdk.NewAttribute("spot_price", res.SpotPrice.String()),
		sdk.NewAttribute("swap_fee", res.SwapFee.String()),
		sdk.NewAttribute("discount", res.Discount.String()),
		sdk.NewAttribute("recipient", res.Recipient),
	)
}

func NewExecuteMarketBuySpotOrderEvt(order SpotOrder, res *ammtypes.MsgSwapByDenomResponse) sdk.Event {
	// convert order price to json string
	orderPrice, err := json.Marshal(order.OrderPrice)
	if err != nil {
		panic(err)
	}

	return sdk.NewEvent(TypeEvtExecuteMarketBuySpotOrder,
		sdk.NewAttribute("order_type", order.OrderType.String()),
		sdk.NewAttribute("order_id", strconv.FormatUint(order.OrderId, 10)),
		sdk.NewAttribute("order_price", string(orderPrice)),
		sdk.NewAttribute("order_amount", order.OrderAmount.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_target_denom", order.OrderTargetDenom),
		sdk.NewAttribute("date", order.Date.String()),
		sdk.NewAttribute("amount", res.Amount.String()),
		sdk.NewAttribute("spot_price", res.SpotPrice.String()),
		sdk.NewAttribute("swap_fee", res.SwapFee.String()),
		sdk.NewAttribute("discount", res.Discount.String()),
		sdk.NewAttribute("recipient", res.Recipient),
	)
}

func NewExecuteLimitOpenPerpetualOrderEvt(order PerpetualOrder, positionId uint64) sdk.Event {
	// convert trigger price to json string
	triggerPrice, err := json.Marshal(order.TriggerPrice)
	if err != nil {
		panic(err)
	}

	return sdk.NewEvent(TypeEvtExecuteLimitOpenPerpetualOrder,
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("position", order.Position.String()),
		sdk.NewAttribute("position_id", strconv.FormatInt(int64(positionId), 10)),
		sdk.NewAttribute("trigger_price", string(triggerPrice)),
	)
}

func NewDeletePendingPerpetualOrderEvt(order PerpetualOrder) sdk.Event {
	return sdk.NewEvent(TypeEvtDeletePendingPerpetualOrder,
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("position", order.Position.String()),
	)
}

func NewExecuteLimitClosePerpetualOrderEvt(order PerpetualOrder, closeAmount string) sdk.Event {
	return sdk.NewEvent(TypeEvtExecuteLimitClosePerpetualOrder,
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("position", order.Position.String()),
		sdk.NewAttribute("close_amount", closeAmount),
	)
}

func NewCreatePerpetualOpenOrderEvt(order PerpetualOrder) sdk.Event {
	return sdk.NewEvent(TypeEvtCreatePerpetualOpenOrder,
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("position", order.Position.String()),
	)
}

func NewCreatePerpetualCloseOrderEvt(order PerpetualOrder) sdk.Event {
	return sdk.NewEvent(TypeEvtCreatePerpetualCloseOrder,
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
	)
}

func NewUpdatePerpetualOrderEvt(order PerpetualOrder, triggerPrice string) sdk.Event {
	return sdk.NewEvent(TypeEvtUpdatePerpetualOrder,
		sdk.NewAttribute("order_id", strconv.FormatInt(int64(order.OrderId), 10)),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("order_type", order.PerpetualOrderType.String()),
		sdk.NewAttribute("position", order.Position.String()),
		sdk.NewAttribute("trigger_price", triggerPrice),
	)
}
