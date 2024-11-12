package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

func NewMsgRemoveAssetInfo(authority, denom string) *MsgRemoveAssetInfo {
	return &MsgRemoveAssetInfo{
		Authority: authority,
		Denom:     denom,
	}
}

// Implements Msg Interface
var _ sdk.Msg = &MsgRemoveAssetInfo{}

func (msg *MsgRemoveAssetInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

// NewMsgAddPriceFeeders creates a new MsgAddPriceFeeders instance
func NewMsgAddPriceFeeders(
	authority string,
	feeders []string,
) *MsgAddPriceFeeders {
	return &MsgAddPriceFeeders{
		Authority: authority,
		Feeders:   feeders,
	}
}

func (msg *MsgAddPriceFeeders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	for _, feeder := range msg.Feeders {
		_, err := sdk.AccAddressFromBech32(feeder)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
		}
	}

	return nil
}

func NewMsgRemovePriceFeeders(authority string, feeders []string) *MsgRemovePriceFeeders {
	return &MsgRemovePriceFeeders{
		Authority: authority,
		Feeders:   feeders,
	}
}

// Implements Msg Interface
var _ sdk.Msg = &MsgRemovePriceFeeders{}

func (msg *MsgRemovePriceFeeders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}

	for _, feeder := range msg.Feeders {
		_, err := sdk.AccAddressFromBech32(feeder)
		if err != nil {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
		}
	}

	return nil
}
