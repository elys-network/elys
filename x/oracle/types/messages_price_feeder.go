package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgSetPriceFeeder    = "update_price_feeder"
	TypeMsgDeletePriceFeeder = "delete_price_feeder"
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

func (msg *MsgSetPriceFeeder) Route() string {
	return RouterKey
}

func (msg *MsgSetPriceFeeder) Type() string {
	return TypeMsgSetPriceFeeder
}

func (msg *MsgSetPriceFeeder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetPriceFeeder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetPriceFeeder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
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
func (msg *MsgDeletePriceFeeder) Route() string {
	return RouterKey
}

func (msg *MsgDeletePriceFeeder) Type() string {
	return TypeMsgDeletePriceFeeder
}

func (msg *MsgDeletePriceFeeder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeletePriceFeeder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeletePriceFeeder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Feeder)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid feeder address (%s)", err)
	}
	return nil
}
