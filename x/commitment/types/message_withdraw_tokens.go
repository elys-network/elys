package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgClaimReward{}

func NewMsgClaimReward(creator string, amount math.Int, denom string) *MsgClaimReward {
	return &MsgClaimReward{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgClaimReward) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if msg.Amount.IsNil() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be nil")
	}

	if msg.Amount.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount cannot be negative")
	}

	return nil
}
