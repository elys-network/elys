package types

import (
	errorsmod "cosmossdk.io/errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateParams{}
var _ sdk.Msg = &MsgAddPriceFeeders{}
var _ sdk.Msg = &MsgRemoveAssetInfo{}
var _ sdk.Msg = &MsgRemovePriceFeeders{}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if err = msg.Params.Validate(); err != nil {
		return err
	}
	return nil
}

func (msg *MsgRemoveAssetInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if err = sdk.ValidateDenom(msg.Denom); err != nil {
		return err
	}
	return nil
}

func (msg *MsgAddPriceFeeders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if len(msg.Feeders) == 0 {
		return fmt.Errorf("no feeders specified")
	}

	for _, feeder := range msg.Feeders {
		_, err = sdk.AccAddressFromBech32(feeder)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
		}
	}

	return nil
}

func (msg *MsgRemovePriceFeeders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	if len(msg.Feeders) == 0 {
		return fmt.Errorf("no feeders specified")
	}
	for _, feeder := range msg.Feeders {
		_, err = sdk.AccAddressFromBech32(feeder)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
		}
	}

	return nil
}
