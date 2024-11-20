package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdatePoolMultipliers{}

func (msg *MsgUpdatePoolMultipliers) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if len(msg.PoolMultipliers) == 0 {
		return errorsmod.Wrapf(ErrInvalidPoolMultiplier, "pool multipliers is empty")
	}

	for _, multiplier := range msg.PoolMultipliers {
		if multiplier.Multiplier.IsNil() {
			return errorsmod.Wrapf(ErrInvalidPoolMultiplier, "multiplier is empty")
		}
		if multiplier.Multiplier.IsNegative() {
			return errorsmod.Wrapf(ErrInvalidPoolMultiplier, "multiplier cannot be negative")
		}
	}

	return nil
}
