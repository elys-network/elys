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
	if msg.InitialMarginRatio.IsNil() || msg.InitialMarginRatio.IsNegative() || msg.InitialMarginRatio.IsZero() {
		return fmt.Errorf("initial margin ratio cannot be negative or zero")
	}
	if msg.MaintenanceMarginRatio.IsNil() || msg.MaintenanceMarginRatio.IsNegative() || msg.MaintenanceMarginRatio.IsZero() {
		return fmt.Errorf("maintenance margin ratio cannot be negative or zero")
	}
	if msg.MaintenanceMarginRatio.GTE(msg.InitialMarginRatio) {
		return fmt.Errorf("maintenance margin ratio cannot be greater than or equal to initial margin ratio")
	}
	imrMMRRatio := msg.InitialMarginRatio.Quo(msg.MaintenanceMarginRatio)
	if imrMMRRatio.LT(MinimumIMRByMMR) || imrMMRRatio.GT(MaximumIMRByMMR) {
		return fmt.Errorf("imr/mmr must be between %s and %s", MinimumIMRByMMR.String(), MaximumIMRByMMR.String())
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
