package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClose = "close"

var _ sdk.Msg = &MsgClose{}

func NewMsgClose(creator string, id uint64) *MsgClose {
	return &MsgClose{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgClose) Route() string {
	return RouterKey
}

func (msg *MsgClose) Type() string {
	return TypeMsgClose
}

func (msg *MsgClose) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClose) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
