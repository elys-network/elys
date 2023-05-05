package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

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
