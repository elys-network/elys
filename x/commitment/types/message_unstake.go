package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	paramtypes "github.com/elys-network/elys/v5/x/parameter/types"
)

var _ sdk.Msg = &MsgUnstake{}

func NewMsgUnstake(creator string, amount math.Int, asset string, validatorAddress string) *MsgUnstake {
	return &MsgUnstake{
		Creator:          creator,
		Amount:           amount,
		Asset:            asset,
		ValidatorAddress: validatorAddress,
	}
}

func (msg *MsgUnstake) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Amount.IsNil() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be nil")
	}

	if msg.Amount.IsNegative() || msg.Amount.IsZero() {
		return errorsmod.Wrapf(ErrInvalidAmount, "Amount can not be negative or zero")
	}

	if err = sdk.ValidateDenom(msg.Asset); err != nil {
		return errorsmod.Wrapf(ErrInvalidDenom, msg.Asset)
	}

	if msg.Asset == paramtypes.Elys {
		_, err = sdk.ValAddressFromBech32(msg.ValidatorAddress)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid validator address  (%s)", err)
		}
	}

	return nil
}
