package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCommitUnclaimedRewards = "commit_tokens"

var _ sdk.Msg = &MsgCommitUnclaimedRewards{}

func NewMsgCommitUnclaimedRewards(creator string, amount sdk.Int, denom string) *MsgCommitUnclaimedRewards {
	return &MsgCommitUnclaimedRewards{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgCommitUnclaimedRewards) Route() string {
	return RouterKey
}

func (msg *MsgCommitUnclaimedRewards) Type() string {
	return TypeMsgCommitUnclaimedRewards
}

func (msg *MsgCommitUnclaimedRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCommitUnclaimedRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCommitUnclaimedRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return sdkerrors.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
