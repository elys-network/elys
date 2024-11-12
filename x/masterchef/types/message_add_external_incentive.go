package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgAddExternalIncentive{}

func (msg *MsgAddExternalIncentive) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.FromBlock >= msg.ToBlock {
		return ErrInvalidBlockRange
	}

	if msg.AmountPerBlock.IsNil() {
		return errorsmod.Wrapf(ErrInvalidAmountPerBlock, "amount per block is nil")
	}

	if msg.AmountPerBlock.IsZero() {
		return errorsmod.Wrapf(ErrInvalidAmountPerBlock, "amount per block is zero")
	}

	if msg.AmountPerBlock.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidAmountPerBlock, "amount per block is negative")
	}

	return nil
}
