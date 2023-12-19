package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimRewards = "claim_rewards"

var _ sdk.Msg = &MsgClaimRewards{}

func NewMsgClaimRewards(creator string) *MsgClaimRewards {
	return &MsgClaimRewards{
		Creator: creator,
	}
}

func (msg *MsgClaimRewards) Route() string {
	return RouterKey
}

func (msg *MsgClaimRewards) Type() string {
	return TypeMsgClaimRewards
}

func (msg *MsgClaimRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
