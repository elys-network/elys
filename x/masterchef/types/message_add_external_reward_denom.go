package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgAddExternalRewardDenom = "add_external_reward_denom"

var _ sdk.Msg = &MsgAddExternalRewardDenom{}

func (msg *MsgAddExternalRewardDenom) Route() string {
	return RouterKey
}

func (msg *MsgAddExternalRewardDenom) Type() string {
	return TypeMsgAddExternalRewardDenom
}

func (msg *MsgAddExternalRewardDenom) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}

func (msg *MsgAddExternalRewardDenom) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddExternalRewardDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return nil
}
