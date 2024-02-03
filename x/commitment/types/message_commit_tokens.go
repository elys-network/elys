package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCommitClaimedRewards = "commit_tokens"

var _ sdk.Msg = &MsgCommitClaimedRewards{}

func NewMsgCommitClaimedRewards(creator string, amount math.Int, denom string) *MsgCommitClaimedRewards {
	return &MsgCommitClaimedRewards{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgCommitClaimedRewards) Route() string {
	return RouterKey
}

func (msg *MsgCommitClaimedRewards) Type() string {
	return TypeMsgCommitClaimedRewards
}

func (msg *MsgCommitClaimedRewards) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCommitClaimedRewards) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCommitClaimedRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
