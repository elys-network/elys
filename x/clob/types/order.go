package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	oracletypes "github.com/elys-network/elys/x/oracle/types"
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

func (order PerpetualOrder) GetQuantity(assetInfo oracletypes.AssetInfo, tradePrice math.LegacyDec) math.LegacyDec {
	denomPrice := tradePrice.Quo(math.LegacyNewDec(10).Power(assetInfo.Decimal))
	return order.Collateral.ToLegacyDec().Mul(order.Leverage).Quo(denomPrice)
}

func (order PerpetualOrder) GetOwnerAccAddress() sdk.AccAddress {
	return sdk.MustAccAddressFromBech32(order.Owner)
}
