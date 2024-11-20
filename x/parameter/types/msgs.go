package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateMinCommission{}

func (msg *MsgUpdateMinCommission) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.MinCommission.IsNil() {
		return fmt.Errorf("min commission cannot be nil")
	}

	if msg.MinCommission.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("min commission cannot be <= 0")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateMinCommission{}

func (msg *MsgUpdateMaxVotingPower) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.MaxVotingPower.IsNil() {
		return fmt.Errorf("max voting power cannot be nil")
	}

	if msg.MaxVotingPower.LTE(math.LegacyZeroDec()) {
		return fmt.Errorf("max voting power cannot be <= 0")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateMinCommission{}

func (msg *MsgUpdateMinSelfDelegation) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.MinSelfDelegation.IsNil() {
		return fmt.Errorf("min self delegation cannot be nil")
	}

	if msg.MinSelfDelegation.LTE(math.ZeroInt()) {
		return fmt.Errorf("min self delegation cannot be <= 0")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateTotalBlocksPerYear{}

func (msg *MsgUpdateTotalBlocksPerYear) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	if msg.TotalBlocksPerYear == 0 {
		return fmt.Errorf("invalid total blocks per year")
	}

	return nil
}
