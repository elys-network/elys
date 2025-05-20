package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ptypes "github.com/elys-network/elys/v4/x/parameter/types"
)

var _ sdk.Msg = &MsgCancelVest{}

func NewMsgCancelVest(creator string, amount math.Int, denom string) *MsgCancelVest {
	return &MsgCancelVest{
		Creator: creator,
		Amount:  amount,
		Denom:   denom,
	}
}

func (msg *MsgCancelVest) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Denom != ptypes.Eden {
		return errorsmod.Wrapf(ErrInvalidDenom, "denom: %s", msg.Denom)
	}

	if msg.Amount.IsNil() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be nil")
	}

	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be negative or zero")
	}

	return nil
}
