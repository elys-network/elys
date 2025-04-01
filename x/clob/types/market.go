package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

func (market PerpetualMarket) ValidateMsgOpenPosition(msg MsgPlaceLimitOrder) error {
	if market.Id != msg.MarketId {
		return ErrPerpetualMarketNotFound
	}
	//if msg.Quantity.LT(market.MinNotional) {
	//	return ErrAmountLessThanNotional
	//}
	//division := msg.Quantity.ToDec().QuoInt(market.MinQuantityTickSize)
	//if !division.TruncateInt().ToDec().Equal(division) {
	//	return fmt.Errorf("amount should in multiples of %s", market.MinQuantityTickSize.String())
	//}
	////division = msg.TriggerPrice.Quo(market.MinPriceTickSize)
	////if !division.TruncateInt().ToDec().Equal(division) {
	////	return fmt.Errorf("price should in multiples of %s", market.MinPriceTickSize.String())
	////}
	//
	//maxLeverage := utils.OneDec.Quo(market.InitialMarginRatio)
	//if msg.Leverage.ToDec().GT(maxLeverage) {
	//	return fmt.Errorf("leverage should not be greater than %s", maxLeverage.String())
	//}

	return nil
}

func (market PerpetualMarket) GetAccount() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/perpetual/%d", market.Id))
}

func (market PerpetualMarket) GetTreasuryAccount() sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("clob/perpetual/treasury/%d", market.Id))
}
