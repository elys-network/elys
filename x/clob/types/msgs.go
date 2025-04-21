package types

import (
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TODO Validate Basic

var _ sdk.Msg = &MsgCreatPerpetualMarket{}

func (msg MsgCreatPerpetualMarket) ValidateBasic() (err error) {
	err = sdk.ValidateDenom(msg.BaseDenom)
	if err != nil {
		return err
	}
	err = sdk.ValidateDenom(msg.QuoteDenom)
	if err != nil {
		return err
	}
	if msg.TwapPricesWindow <= 10 {
		return fmt.Errorf("max twap prices time must be greater than 10")
	}
	if msg.MaxAbsFundingRate.IsNil() || msg.MaxAbsFundingRate.IsNegative() {
		return fmt.Errorf("max abs funding rate cannot be negative or nil")
	}
	if msg.MaxAbsFundingRateChange.IsNil() || msg.MaxAbsFundingRateChange.IsNegative() {
		return fmt.Errorf("max abs funding rate cannot be negative or nil")
	}
	return nil
}

func (msg MsgDeposit) ValidateBasic() error {
	return nil
}

func (msg MsgPlaceLimitOrder) ValidateBasic() error {
	return nil
}

func (msg MsgPlaceMarketOrder) ValidateBasic() error {
	return nil
}

func (msg MsgUpdateParams) ValidateBasic() error {
	return nil
}
