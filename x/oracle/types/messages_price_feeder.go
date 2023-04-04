package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePriceFeeder = "create_price_feeder"
	TypeMsgUpdatePriceFeeder = "update_price_feeder"
	TypeMsgDeletePriceFeeder = "delete_price_feeder"
)

var _ sdk.Msg = &MsgCreatePriceFeeder{}

func NewMsgCreatePriceFeeder(
	creator string,
	feeder string,
	isActive bool,
) *MsgCreatePriceFeeder {
	return &MsgCreatePriceFeeder{
		Creator:  creator,
		Feeder:   feeder,
		IsActive: isActive,
	}
}

func (msg *MsgCreatePriceFeeder) Route() string {
	return RouterKey
}

func (msg *MsgCreatePriceFeeder) Type() string {
	return TypeMsgCreatePriceFeeder
}

func (msg *MsgCreatePriceFeeder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePriceFeeder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePriceFeeder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePriceFeeder{}

func NewMsgUpdatePriceFeeder(
	creator string,
	feeder string,
	isActive bool,
) *MsgUpdatePriceFeeder {
	return &MsgUpdatePriceFeeder{
		Creator:  creator,
		Feeder:   feeder,
		IsActive: isActive,
	}
}

func (msg *MsgUpdatePriceFeeder) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePriceFeeder) Type() string {
	return TypeMsgUpdatePriceFeeder
}

func (msg *MsgUpdatePriceFeeder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePriceFeeder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePriceFeeder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePriceFeeder{}

func NewMsgDeletePriceFeeder(
	creator string,
	feeder string,
) *MsgDeletePriceFeeder {
	return &MsgDeletePriceFeeder{
		Creator: creator,
		Feeder:  feeder,
	}
}
func (msg *MsgDeletePriceFeeder) Route() string {
	return RouterKey
}

func (msg *MsgDeletePriceFeeder) Type() string {
	return TypeMsgDeletePriceFeeder
}

func (msg *MsgDeletePriceFeeder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
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
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
