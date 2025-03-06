package types

import sdk "github.com/cosmos/cosmos-sdk/types"

// TODO Validate Basic

var _ sdk.Msg = &MsgCreatPerpetualMarket{}

func (msg MsgCreatPerpetualMarket) ValidateBasic() error {
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
