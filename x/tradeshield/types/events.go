package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	
)

const (
	TypeEvtCloseOrder   = "close_order"
	
)

func EmitCloseOrdersEvent(ctx sdk.Context, order SpotOrder) {
	ctx.EventManager().EmitEvents(sdk.Events{
		NewCloseOrdersEvt(order),
	})
}

func NewCloseOrdersEvt(order SpotOrder) sdk.Event{
	return sdk.NewEvent(TypeEvtCloseOrder,
		sdk.NewAttribute("order_type", order.OrderType.String()),
		sdk.NewAttribute("owner_address", order.OwnerAddress),
		sdk.NewAttribute("id", strconv.FormatInt(int64(order.OrderId), 10)),
	)
}