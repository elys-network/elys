package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateMinCommission{}

func NewMsgExitPool(creator string, minCommission string) *MsgUpdateMinCommission {
	return &MsgUpdateMinCommission{
		Creator:       creator,
		MinCommission: minCommission,
	}
}

func (msg *MsgUpdateMinCommission) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdateMinCommission) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateMinCommission{}

func NewMsgUpdateMaxVotingPower(creator string, maxVotingPower string) *MsgUpdateMaxVotingPower {
	return &MsgUpdateMaxVotingPower{
		Creator:        creator,
		MaxVotingPower: maxVotingPower,
	}
}

func (msg *MsgUpdateMaxVotingPower) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdateMaxVotingPower) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateMinCommission{}

func NewMsgUpdateMinSelfDelegation(creator string, minSelfDelegation string) *MsgUpdateMinSelfDelegation {
	return &MsgUpdateMinSelfDelegation{
		Creator:           creator,
		MinSelfDelegation: minSelfDelegation,
	}
}

func (msg *MsgUpdateMinSelfDelegation) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdateMinSelfDelegation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}
