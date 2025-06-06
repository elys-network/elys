package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUncommitTokens{}

func NewMsgUncommitTokens(creator string, amount math.Int, denom string) *MsgUncommitTokens {
	return &MsgUncommitTokens{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg MsgUncommitTokens) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "Invalid creator address: %v", err)
	}

	if err = sdk.ValidateDenom(msg.Denom); err != nil {
		return errorsmod.Wrapf(ErrInvalidDenom, msg.Denom)
	}

	if msg.Amount.IsNil() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be nil")
	}

	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be negative or zero")
	}
	return nil
}
