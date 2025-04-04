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
	if msg.MaxTwapPricesTime <= 10 {
		return fmt.Errorf("max twap prices time must be greater than 10")
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
