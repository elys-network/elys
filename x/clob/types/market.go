package types

import (
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (market PerpetualMarket) ValidateOpenPositionRequest(marketId uint64, price, quantity math.LegacyDec, isMarketOrder bool) error {
	if market.Id != marketId {
		return ErrPerpetualMarketNotFound
	}
	if price.Mul(quantity).LT(market.MinNotional) {
		return errors.New("trade value less than minimum notional value")
	}
	if quantity.LT(market.MinQuantityTickSize) {
		return errors.New("quantity less than minimum quantity tick size")
	}
	if !quantity.Quo(market.MinQuantityTickSize).IsInteger() {
		return errors.New("quantity is not of proper tick size")
	}
	if !price.Quo(market.MinPriceTickSize).IsInteger() {
		return errors.New("price is not of proper tick size")
	}
	return nil
}

func (market PerpetualMarket) GetAccount() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/perpetual/%d", market.Id))
}

func (market *PerpetualMarket) UpdateTotalOpenInterest(buyerBefore, sellerBefore, tradeSize math.LegacyDec) {
	if tradeSize.LTE(math.LegacyZeroDec()) {
		panic("trade size cannot be 0 or negative")
	}
	oiBefore := buyerBefore.Abs().Add(sellerBefore.Abs())
	buyerAfter := buyerBefore.Add(tradeSize)
	sellerAfter := sellerBefore.Sub(tradeSize)
	oiAfter := buyerAfter.Abs().Add(sellerAfter.Abs())

	netChange := oiAfter.Sub(oiBefore).QuoInt64(2)
	market.TotalOpen = market.TotalOpen.Add(netChange)
}
