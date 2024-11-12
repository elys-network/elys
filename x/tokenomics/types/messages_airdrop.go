package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

func (msg *MsgUpdateAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate Intent is not empty
	if len(msg.Intent) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "intent cannot be empty")
	}

	// Validate Amount is positive
	if msg.Amount <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidCoins, "amount must be positive")
	}

	// Validate Expiry
	if msg.Expiry <= 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "expiry must be a positive timestamp")
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

func (msg *MsgDeleteAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	// Validate Intent is not empty
	if len(msg.Intent) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrUnknownRequest, "intent cannot be empty")
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

func (msg *MsgClaimAirdrop) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
