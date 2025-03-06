package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (order PerpetualOrder) IsBuy() bool {
	return IsBuy(order.OrderType)
}

func IsBuy(orderType OrderType) bool {
	switch orderType {
	case OrderType_ORDER_TYPE_LIMIT_BUY, OrderType_ORDER_TYPE_MARKET_BUY:
		return true
	default:
		return false
	}
}

func (order PerpetualOrder) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(order.Owner)
}
