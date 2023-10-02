package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUnbond = "unbond"

var _ sdk.Msg = &MsgUnbond{}

func NewMsgUnbond(creator string, amount string) *MsgUnbond {
	return &MsgUnbond{
		Creator: creator,
		Amount:  amount,
	}
}

func (msg *MsgUnbond) Route() string {
	return RouterKey
}

func (msg *MsgUnbond) Type() string {
	return TypeMsgUnbond
}

func (msg *MsgUnbond) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUnbond) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUnbond) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
