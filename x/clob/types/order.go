package types

import (
	sdkmath "cosmossdk.io/math"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

const PriceMultiplier int64 = 1_000_000

func NewOrderId(marketId uint64, orderType OrderType, priceTick int64, counter uint64) OrderId {
	return OrderId{
		MarketId:  marketId,
		OrderType: orderType,
		PriceTick: priceTick,
		Counter:   counter,
	}
}

func (o OrderId) KeyWithoutPrefix() []byte {
	key := sdk.Uint64ToBigEndian(o.MarketId)
	key = append(key, []byte("/")...)
	orderTypeByte := FalseByte
	counterBytes := sdk.Uint64ToBigEndian(o.Counter)
	if IsBuy(o.OrderType) {
		orderTypeByte = TrueByte
		counterBytes = sdk.Uint64ToBigEndian(math.MaxUint64 - o.Counter) // Subtracting it so that in buy order book, it's sorted by counter (if 2 orders have same price) as Reverse iterator will be used
	}
	key = append(key, orderTypeByte)
	key = append(key, []byte("/")...)
	priceBytes := sdk.Uint64ToBigEndian(uint64(o.PriceTick))
	key = append(key, priceBytes...)
	key = append(key, []byte("/")...)
	key = append(key, counterBytes...)
	return key
}

func NewPerpetualOrder(marketId uint64, orderType OrderType, price sdkmath.LegacyDec, counter uint64, owner sdk.AccAddress, amount, filled sdkmath.LegacyDec, subAccountId uint64) PerpetualOrder {
	orderId := NewOrderId(marketId, orderType, price.MulInt64(PriceMultiplier).TruncateInt64(), counter)
	return PerpetualOrder{
		OrderId:      orderId,
		Owner:        owner.String(),
		Amount:       amount,
		Filled:       filled,
		SubAccountId: subAccountId,
	}
}

func (order PerpetualOrder) GetMarketId() uint64 {
	return order.OrderId.MarketId
}

func (order PerpetualOrder) GetOrderType() OrderType {
	return order.OrderId.OrderType
}

func (order PerpetualOrder) GetPriceTick() int64 {
	return order.OrderId.PriceTick
}

func (order PerpetualOrder) GetCounter() uint64 {
	return order.OrderId.Counter
}

func (order PerpetualOrder) IsBuy() bool {
	return IsBuy(order.GetOrderType())
}

func IsBuy(orderType OrderType) bool {
	switch orderType {
	case OrderType_ORDER_TYPE_LIMIT_BUY, OrderType_ORDER_TYPE_MARKET_BUY:
		return true
	default:
		return false
	}
}

func (order PerpetualOrder) GetPrice() sdkmath.LegacyDec {
	return sdkmath.LegacyNewDec(order.OrderId.PriceTick).QuoInt64(PriceMultiplier)
}

func (order PerpetualOrder) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(order.Owner)
}

func (order PerpetualOrder) UnfilledValue() sdkmath.LegacyDec {
	return order.GetPrice().Mul(order.Amount.Sub(order.Filled))
}
