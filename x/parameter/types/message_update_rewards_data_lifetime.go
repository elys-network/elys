package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgUpdateRewardsDataLifetime = "update_rewards_data_lifetime"

var _ sdk.Msg = &MsgUpdateRewardsDataLifetime{}

func NewMsgUpdateRewardsDataLifetime(creator string, rewardsDataLifetime string) *MsgUpdateRewardsDataLifetime {
	return &MsgUpdateRewardsDataLifetime{
		Creator:             creator,
		RewardsDataLifetime: rewardsDataLifetime,
	}
}

func (msg *MsgUpdateRewardsDataLifetime) Route() string {
	return RouterKey
}

func (msg *MsgUpdateRewardsDataLifetime) Type() string {
	return TypeMsgUpdateRewardsDataLifetime
}

func (msg *MsgUpdateRewardsDataLifetime) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateRewardsDataLifetime) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateRewardsDataLifetime) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	res, ok := sdk.NewIntFromString(msg.RewardsDataLifetime)
	if !ok || !res.IsPositive() {
		return sdkerrors.Wrapf(ErrInvalidRewardsDataLifecycle, "invalid data in rewards_data_lifecycle")
	}
	return nil
}
