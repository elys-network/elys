package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawTokens = "withdraw_tokens"

var _ sdk.Msg = &MsgWithdrawTokens{}

func NewMsgWithdrawTokens(creator string, amount sdk.Int, denom string) *MsgWithdrawTokens {
	return &MsgWithdrawTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgWithdrawTokens) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawTokens) Type() string {
	return TypeMsgWithdrawTokens
}

func (msg *MsgWithdrawTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
