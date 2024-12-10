package types

import (
	"encoding/json"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
)

const (
	TypeEvtCloseSpotOrder                 = "tradeshield/close_spot_order"
	TypeEvtCancelPerpetualOrder           = "tradeshield/cancel_perpetual_order"
	TypeEvtExecuteOrders                  = "tradeshield/execute_orders"
	TypeEvtExecuteLimitOpenPerpetualOrder = "tradeshield/execute_perpetual_limit_open_order"
	TypeEvtExecuteSpotOrder               = "tradeshield/execute_spot_order"
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
		sdk.NewAttribute("trading_asset", order.TradingAsset),
		sdk.NewAttribute("leverage", order.Leverage.String()),
		sdk.NewAttribute("take_profit_price", order.TakeProfitPrice.String()),
		sdk.NewAttribute("position_id", strconv.FormatInt(int64(order.PositionId), 10)),
		sdk.NewAttribute("status", order.Status.String()),
		sdk.NewAttribute("stop_loss_price", order.StopLossPrice.String()),
	)
}

func NewExecuteSpotOrderEvt(order SpotOrder, res *ammtypes.MsgSwapByDenomResponse) sdk.Event {
	// convert order price to json string
	orderPrice, err := json.Marshal(order.OrderPrice)
	if err != nil {
		panic(err)
	}

	return sdk.NewEvent(TypeEvtExecuteSpotOrder,
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
