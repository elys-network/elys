package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUncommitTokens = "uncommit_tokens"

var _ sdk.Msg = &MsgUncommitTokens{}

func NewMsgUncommitTokens(creator string, amount sdk.Int, denom string) *MsgUncommitTokens {
	return &MsgUncommitTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgUncommitTokens) Route() string {
	return RouterKey
}

func (msg *MsgUncommitTokens) Type() string {
	return TypeMsgUncommitTokens
}

func (msg *MsgUncommitTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUncommitTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUncommitTokens) ValidateBasic() error { // TODO
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
