package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgClaimReward = "withdraw_tokens"

var _ sdk.Msg = &MsgClaimReward{}

func NewMsgClaimReward(creator string, amount sdk.Int, denom string) *MsgClaimReward {
	return &MsgClaimReward{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgClaimReward) Route() string {
	return RouterKey
}

func (msg *MsgClaimReward) Type() string {
	return TypeMsgClaimReward
}

func (msg *MsgClaimReward) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClaimReward) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClaimReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
