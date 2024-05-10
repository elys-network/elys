package types

import (
	errorsmod "cosmossdk.io/errors"
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateBrokerAddress{}

func NewMsgUpdateBrokerAddress(creator string, brokerAddress string) *MsgUpdateBrokerAddress {
	return &MsgUpdateBrokerAddress{
		Creator:       creator,
		BrokerAddress: brokerAddress,
	}
}

func (msg *MsgUpdateBrokerAddress) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdateBrokerAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	_, err = sdk.AccAddressFromBech32(msg.BrokerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid broker address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateTotalBlocksPerYear{}

func NewMsgUpdateTotalBlocksPerYear(creator string, totalBlocksPerYear int64) *MsgUpdateTotalBlocksPerYear {
	return &MsgUpdateTotalBlocksPerYear{
		Creator:            creator,
		TotalBlocksPerYear: totalBlocksPerYear,
	}
}

func (msg *MsgUpdateTotalBlocksPerYear) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgUpdateTotalBlocksPerYear) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.TotalBlocksPerYear == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid total blocks per year")
	}

	return nil
}
