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

	return nil
}
