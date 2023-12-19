package types

import (
	errorsmod "cosmossdk.io/errors"
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

func (msg MsgUncommitTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
