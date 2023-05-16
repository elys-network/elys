package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgSetWithdrawAddress = "set_withdraw_address"

var _ sdk.Msg = &MsgSetWithdrawAddress{}

func NewMsgSetWithdrawAddress(delegatorAddress string, withdrawAddress string) *MsgSetWithdrawAddress {
	return &MsgSetWithdrawAddress{
		DelegatorAddress: delegatorAddress,
		WithdrawAddress:  withdrawAddress,
	}
}

func (msg *MsgSetWithdrawAddress) Route() string {
	return RouterKey
}

func (msg *MsgSetWithdrawAddress) Type() string {
	return TypeMsgSetWithdrawAddress
}

func (msg *MsgSetWithdrawAddress) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetWithdrawAddress) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetWithdrawAddress) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

const TypeMsgWithdrawDelegatorReward = "withdraw_delegator_reward"

var _ sdk.Msg = &MsgWithdrawDelegatorReward{}

func NewMsgWithdrawDelegatorReward(creator string, validatorAddress string) *MsgWithdrawDelegatorReward {
	return &MsgWithdrawDelegatorReward{
		DelegatorAddress: creator,
		ValidatorAddress: validatorAddress,
	}
}

func (msg *MsgWithdrawDelegatorReward) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawDelegatorReward) Type() string {
	return TypeMsgWithdrawDelegatorReward
}

func (msg *MsgWithdrawDelegatorReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawDelegatorReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawDelegatorReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

const TypeMsgWithdrawValidatorCommission = "withdraw_validator_commission"

var _ sdk.Msg = &MsgWithdrawValidatorCommission{}

func NewMsgWithdrawValidatorCommission(creator string) *MsgWithdrawValidatorCommission {
	return &MsgWithdrawValidatorCommission{
		ValidatorAddress: creator,
	}
}

func (msg *MsgWithdrawValidatorCommission) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawValidatorCommission) Type() string {
	return TypeMsgWithdrawValidatorCommission
}

func (msg *MsgWithdrawValidatorCommission) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawValidatorCommission) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.ValidatorAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
