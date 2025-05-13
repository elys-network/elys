package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func NewPerpetualOrder(marketId uint64, orderType OrderType, price math.LegacyDec, height uint64, owner sdk.AccAddress, amount, filled math.LegacyDec) PerpetualOrder {
	return PerpetualOrder{
		MarketId:    marketId,
		OrderType:   orderType,
		Price:       price,
		BlockHeight: height,
		Owner:       owner.String(),
		Amount:      amount,
		Filled:      filled,
	}
}

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
