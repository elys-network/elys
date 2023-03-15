package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgDepositTokens = "deposit_tokens"

var _ sdk.Msg = &MsgDepositTokens{}

func NewMsgDepositTokens(creator string, amount sdk.Int, denom string) *MsgDepositTokens {
	return &MsgDepositTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgDepositTokens) Route() string {
	return RouterKey
}

func (msg *MsgDepositTokens) Type() string {
	return TypeMsgDepositTokens
}

func (msg *MsgDepositTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDepositTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDepositTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
