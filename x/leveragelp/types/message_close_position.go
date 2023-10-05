package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClosePosition = "close_position"

var _ sdk.Msg = &MsgClosePosition{}

func NewMsgClosePosition(creator string) *MsgClosePosition {
	return &MsgClosePosition{
		Creator: creator,
	}
}

func (msg *MsgClosePosition) Route() string {
	return RouterKey
}

func (msg *MsgClosePosition) Type() string {
	return TypeMsgClosePosition
}

func (msg *MsgClosePosition) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClosePosition) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClosePosition) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
