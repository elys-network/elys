package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgSetPriceFeeder{}

func NewMsgSetPriceFeeder(
	feeder string,
	isActive bool,
) *MsgSetPriceFeeder {
	return &MsgSetPriceFeeder{
		Feeder:   feeder,
		IsActive: isActive,
	}
}

func (msg *MsgSetPriceFeeder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePriceFeeder{}

func NewMsgDeletePriceFeeder(
	feeder string,
) *MsgDeletePriceFeeder {
	return &MsgDeletePriceFeeder{
		Feeder: feeder,
	}
}

func (msg *MsgDeletePriceFeeder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
	}
	return nil
}
