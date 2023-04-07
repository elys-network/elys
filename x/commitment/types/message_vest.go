package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgVest = "vest"

var _ sdk.Msg = &MsgVest{}

func NewMsgVest(creator string, amount sdk.Int, denom string) *MsgVest {
	return &MsgVest{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgVest) Route() string {
	return RouterKey
}

func (msg *MsgVest) Type() string {
	return TypeMsgVest
}

func (msg *MsgVest) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgVest) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgVest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
