package types

import (
	"cosmossdk.io/math"
)

type Trade struct {
	BuyerSubAccount  SubAccount
	SellerSubAccount SubAccount
	MarketId         uint64
	Price            math.LegacyDec
	Quantity         math.LegacyDec
}

func (t Trade) GetRequiredInitialMargin(market PerpetualMarket) math.Int {
	value := t.Price.Mul(t.Quantity)
	requiredInitialMargin := value.Mul(market.InitialMarginRatio).RoundInt()
	return requiredInitialMargin
}

func (t Trade) GetTradeValue() math.LegacyDec {
	return t.Price.Mul(t.Quantity)
}
