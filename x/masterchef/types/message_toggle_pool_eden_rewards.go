package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgTogglePoolEdenRewards = "toggle_pool_eden_rewards"

var _ sdk.Msg = &MsgTogglePoolEdenRewards{}

func (msg *MsgTogglePoolEdenRewards) Route() string {
	return RouterKey
}

func (msg *MsgTogglePoolEdenRewards) Type() string {
	return TypeMsgTogglePoolEdenRewards
}

func (msg *MsgTogglePoolEdenRewards) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgTogglePoolEdenRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgTogglePoolEdenRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
