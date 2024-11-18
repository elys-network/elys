package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateVestingInfo{}

func (msg *MsgUpdateVestingInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	vestingInfo := VestingInfo{
		BaseDenom:      msg.BaseDenom,
		VestingDenom:   msg.VestingDenom,
		NumBlocks:      msg.NumBlocks,
		VestNowFactor:  sdkmath.NewInt(msg.VestNowFactor),
		NumMaxVestings: msg.NumMaxVestings,
	}

	if err = vestingInfo.Validate(); err != nil {
		return err
	}
	return nil
}
