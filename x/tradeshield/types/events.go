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
	)
}