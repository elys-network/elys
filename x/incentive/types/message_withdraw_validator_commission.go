package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
