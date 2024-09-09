package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgWithdrawElysStakingRewards = "withdraw_elys_staking_rewards"

var _ sdk.Msg = &MsgWithdrawElysStakingRewards{}

func (msg *MsgWithdrawElysStakingRewards) Route() string {
	return RouterKey
}

func (msg *MsgWithdrawElysStakingRewards) Type() string {
	return TypeMsgWithdrawElysStakingRewards
}

func (msg *MsgWithdrawElysStakingRewards) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgWithdrawElysStakingRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWithdrawElysStakingRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.DelegatorAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}
