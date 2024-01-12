package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAirdrop = "create_airdrop"
	TypeMsgUpdateAirdrop = "update_airdrop"
	TypeMsgDeleteAirdrop = "delete_airdrop"
	TypeMsgClaimAirdrop  = "claim_airdrop"
)

var _ sdk.Msg = &MsgCreateAirdrop{}

func NewMsgCreateAirdrop(
	authority string,
	intent string,
	amount uint64,
	expiry uint64,
) *MsgCreateAirdrop {
	return &MsgCreateAirdrop{
		Authority: authority,
		Intent:    intent,
		Amount:    amount,
		Expiry:    expiry,
	}
}

func (msg *MsgCreateAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgCreateAirdrop) Type() string {
	return TypeMsgCreateAirdrop
}

func (msg *MsgCreateAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgCreateAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAirdrop{}

func NewMsgUpdateAirdrop(
	authority string,
	intent string,
	amount uint64,
	expiry uint64,
) *MsgUpdateAirdrop {
	return &MsgUpdateAirdrop{
		Authority: authority,
		Intent:    intent,
		Amount:    amount,
		Expiry:    expiry,
	}
}

func (msg *MsgUpdateAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAirdrop) Type() string {
	return TypeMsgUpdateAirdrop
}

func (msg *MsgUpdateAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdateAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteAirdrop{}

func NewMsgDeleteAirdrop(
	authority string,
	intent string,
) *MsgDeleteAirdrop {
	return &MsgDeleteAirdrop{
		Authority: authority,
		Intent:    intent,
	}
}

func (msg *MsgDeleteAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAirdrop) Type() string {
	return TypeMsgDeleteAirdrop
}

func (msg *MsgDeleteAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgDeleteAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgClaimAirdrop{}

func NewMsgClaimAirdrop(
	sender string,
) *MsgClaimAirdrop {
	return &MsgClaimAirdrop{
		Sender: sender,
	}
}

func (msg *MsgClaimAirdrop) Route() string {
	return RouterKey
}

func (msg *MsgClaimAirdrop) Type() string {
	return TypeMsgClaimAirdrop
}

func (msg *MsgClaimAirdrop) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgClaimAirdrop) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
