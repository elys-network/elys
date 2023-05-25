package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// distribution message types
const (
	TypeMsgWithdrawRewards             = "withdraw_reward"
	TypeMsgWithdrawValidatorCommission = "withdraw_validator_commission"
)

// Verify interface at compile time
var _, _ sdk.Msg = &MsgWithdrawRewards{}, &MsgWithdrawValidatorCommission{}

func NewMsgWithdrawRewards(delAddr sdk.AccAddress) *MsgWithdrawRewards {
	return &MsgWithdrawRewards{
		DelegatorAddress: delAddr.String(),
	}
}

func (msg MsgWithdrawRewards) Route() string { return ModuleName }
func (msg MsgWithdrawRewards) Type() string  { return TypeMsgWithdrawRewards }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawRewards) GetSigners() []sdk.AccAddress {
	delegator, _ := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	return []sdk.AccAddress{delegator}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawRewards) ValidateBasic() error {
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	return nil
}

func NewMsgWithdrawValidatorCommission(delAddr sdk.AccAddress, valAddr sdk.ValAddress) *MsgWithdrawValidatorCommission {
	return &MsgWithdrawValidatorCommission{
		DelegatorAddress: delAddr.String(),
		ValidatorAddress: valAddr.String(),
	}
}

func (msg MsgWithdrawValidatorCommission) Route() string { return ModuleName }
func (msg MsgWithdrawValidatorCommission) Type() string  { return TypeMsgWithdrawValidatorCommission }

// Return address that must sign over msg.GetSignBytes()
func (msg MsgWithdrawValidatorCommission) GetSigners() []sdk.AccAddress {
	valAddr, _ := sdk.ValAddressFromBech32(msg.ValidatorAddress)
	return []sdk.AccAddress{sdk.AccAddress(valAddr)}
}

// get the bytes for the message signer to sign on
func (msg MsgWithdrawValidatorCommission) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// quick validity check
func (msg MsgWithdrawValidatorCommission) ValidateBasic() error {
	if _, err := sdk.ValAddressFromBech32(msg.ValidatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid validator address: %s", err)
	}
	if _, err := sdk.AccAddressFromBech32(msg.DelegatorAddress); err != nil {
		return sdkerrors.ErrInvalidAddress.Wrapf("invalid delegator address: %s", err)
	}
	return nil
}
