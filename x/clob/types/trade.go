package types

import (
	"cosmossdk.io/math"
	"github.com/elys-network/elys/utils"
)

type Trade struct {
	BuyerSubAccount, SellerSubAccount SubAccount
	MarketId                          uint64
	Price                             math.Dec
	Quantity                          math.Int
}

func (t Trade) GetRequiredInitialMargin(market PerpetualMarket) (math.Int, error) {
	m1, err := t.Price.Mul(t.QuantityDec())
	if err != nil {
		return math.Int{}, err
	}
	requiredInitialMarginDec, err := m1.Mul(market.InitialMarginRatio)
	if err != nil {
		return math.Int{}, err
	}
	requiredInitialMargin, err := requiredInitialMarginDec.SdkIntTrim()
	if err != nil {
		return math.Int{}, err
	}
	return requiredInitialMargin, nil
}

func (t Trade) QuantityDec() math.Dec {
	return utils.IntToDec(t.Quantity)
}
