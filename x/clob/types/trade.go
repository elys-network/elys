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

func (t Trade) GetTradeValue() math.LegacyDec {
	return t.Price.Mul(t.Quantity)
}

func NewTrade(marketId uint64, qty, price math.LegacyDec, buyer, seller SubAccount) Trade {
	if qty.IsNegative() || qty.IsZero() || qty.IsNil() {
		panic("trade quantity must be positive")
	}
	return Trade{
		MarketId:         marketId,
		Quantity:         qty, // Assumed magnitude
		Price:            price,
		BuyerSubAccount:  buyer,
		SellerSubAccount: seller,
	}
}
