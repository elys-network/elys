package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddEntry = "add_entry"

var _ sdk.Msg = &MsgAddEntry{}

func NewMsgAddEntry(creator string) *MsgAddEntry {
	return &MsgAddEntry{
		Creator: creator,
	}
}

func (msg *MsgAddEntry) Route() string {
	return RouterKey
}

func (msg *MsgAddEntry) Type() string {
	return TypeMsgAddEntry
}

func (msg *MsgAddEntry) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddEntry) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddEntry) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
