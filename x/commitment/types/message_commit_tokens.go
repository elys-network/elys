package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCommitTokens = "commit_tokens"

var _ sdk.Msg = &MsgCommitTokens{}

func NewMsgCommitTokens(creator string, amount sdk.Int, denom string) *MsgCommitTokens {
	return &MsgCommitTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgCommitTokens) Route() string {
	return RouterKey
}

func (msg *MsgCommitTokens) Type() string {
	return TypeMsgCommitTokens
}

func (msg *MsgCommitTokens) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCommitTokens) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCommitTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
