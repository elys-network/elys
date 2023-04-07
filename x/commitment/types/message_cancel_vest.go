package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCancelVest = "cancel_vest"

var _ sdk.Msg = &MsgCancelVest{}

func NewMsgCancelVest(creator string, amount sdk.Int, denom string) *MsgCancelVest {
	return &MsgCancelVest{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgCancelVest) Route() string {
	return RouterKey
}

func (msg *MsgCancelVest) Type() string {
	return TypeMsgCancelVest
}

func (msg *MsgCancelVest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelVest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelVest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
