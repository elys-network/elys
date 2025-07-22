package types

import (
	"errors"

	"cosmossdk.io/math"
)

type Trade struct {
	BuyerSubAccount     SubAccount
	SellerSubAccount    SubAccount
	MarketId            uint64
	Price               math.LegacyDec
	Quantity            math.LegacyDec
	IsBuyerLiquidation  bool
	IsSellerLiquidation bool
	IsBuyerTaker        bool
}

func (t Trade) Validate() error {
	if t.Quantity.IsNil() || t.Quantity.LTE(math.LegacyZeroDec()) {
		return errors.New("trade size cannot be 0 or negative")
	}
	if t.Price.IsNil() || t.Price.LTE(math.LegacyZeroDec()) {
		return errors.New("trade price cannot be 0 or negative")
	}
	if t.MarketId == 0 {
		return errors.New("market id cannot be 0 in a trade")
	}
	return nil
}
func (t Trade) GetTradeValue() math.LegacyDec {
	return t.Price.Mul(t.Quantity)
}

func NewTrade(marketId uint64, qty, price math.LegacyDec, buyer, seller SubAccount) Trade {
	// Note: Input validation is now handled by the Trade.Validate() method
	// This maintains backward compatibility while still providing validation
	return Trade{
		MarketId:         marketId,
		Quantity:         qty, // Assumed magnitude
		Price:            price,
		BuyerSubAccount:  buyer,
		SellerSubAccount: seller,
	}
}
