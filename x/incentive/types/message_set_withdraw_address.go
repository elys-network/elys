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
