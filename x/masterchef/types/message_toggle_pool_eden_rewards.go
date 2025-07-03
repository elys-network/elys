package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgTogglePoolEdenRewards{}

func (msg *MsgTogglePoolEdenRewards) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.PoolId == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "pool id cannot be zero")
	}

	return nil
}
