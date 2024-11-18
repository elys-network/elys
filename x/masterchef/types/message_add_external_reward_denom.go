package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddExternalRewardDenom{}

func (msg *MsgAddExternalRewardDenom) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.MinAmount.IsNil() {
		return errorsmod.Wrapf(ErrInvalidMinAmount, "min amount is nil")
	}

	if msg.MinAmount.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidMinAmount, "min amount is negative")
	}

	if msg.MinAmount.IsZero() {
		return errorsmod.Wrapf(ErrInvalidAmountPerBlock, "min amount is zero")
	}

	if err = sdk.ValidateDenom(msg.RewardDenom); err != nil {
		return err
	}

	return nil
}
