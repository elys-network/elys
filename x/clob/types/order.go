package types

import (
	sdkmath "cosmossdk.io/math"
	"math"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/utils"
)

func NewOrderKey(marketId uint64, orderType OrderType, price sdkmath.LegacyDec, height uint64) OrderKey {
	return OrderKey{
		MarketId:    marketId,
		OrderType:   orderType,
		Price:       price,
		BlockHeight: height,
	}
}

func (o OrderKey) KeyWithoutPrefix() []byte {
	key := sdk.Uint64ToBigEndian(o.MarketId)
	key = append(key, []byte("/")...)
	orderTypeByte := FalseByte
	heightBytes := sdk.Uint64ToBigEndian(o.BlockHeight)
	if IsBuy(o.OrderType) {
		orderTypeByte = TrueByte
		heightBytes = sdk.Uint64ToBigEndian(math.MaxUint64 - o.BlockHeight) // Subtracting it so that in buy order book, it's sorted by height as Reverse iterator will be used
	}
	key = append(key, orderTypeByte)
	key = append(key, []byte("/")...)
	paddedPrice := utils.GetPaddedDecString(o.Price)
	key = append(key, []byte(paddedPrice)...)
	key = append(key, []byte("/")...)
	key = append(key, heightBytes...)
	return key
}

func NewPerpetualOrder(marketId uint64, orderType OrderType, price sdkmath.LegacyDec, height uint64, owner sdk.AccAddress, amount, filled sdkmath.LegacyDec) PerpetualOrder {
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
