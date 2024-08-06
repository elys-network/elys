package types

import (
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgUpdateParams       string = "UpdateParams"
	TypeMsgAddAssetInfo       string = "AddAssetInfo"
	TypeMsgRemoveAssetInfo    string = "RemoveAssetInfo"
	TypeMsgAddPriceFeeders    string = "AddPriceFeeders"
	TypeMsgRemovePriceFeeders string = "RemovePriceFeeders"
)

var _ sdk.Msg = &MsgUpdateParams{}

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleAminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgAddAssetInfo{}

func NewMsgAddAssetInfo(
	authority string,
	denom string,
	display string,
	bandTicker string,
	elysTicker string,
	decimal uint64,
) *MsgAddAssetInfo {
	return &MsgAddAssetInfo{
		Authority:  authority,
		Denom:      denom,
		Display:    display,
		BandTicker: bandTicker,
		ElysTicker: elysTicker,
		Decimal:    decimal,
	}
}

func (msg *MsgAddAssetInfo) Route() string {
	return RouterKey
}

func (msg *MsgAddAssetInfo) Type() string {
	return TypeMsgAddAssetInfo
}

func (msg *MsgAddAssetInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddAssetInfo) GetSignBytes() []byte {
	bz := ModuleAminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddAssetInfo) ValidateBasic() error {
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

func (msg *MsgRemoveAssetInfo) Route() string {
	return RouterKey
}

func (msg *MsgRemoveAssetInfo) Type() string {
	return TypeMsgAddAssetInfo
}

func (msg *MsgRemoveAssetInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemoveAssetInfo) GetSignBytes() []byte {
	bz := ModuleAminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

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

func (msg *MsgAddPriceFeeders) Route() string {
	return RouterKey
}

func (msg *MsgAddPriceFeeders) Type() string {
	return TypeMsgAddAssetInfo
}

func (msg *MsgAddPriceFeeders) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgAddPriceFeeders) GetSignBytes() []byte {
	bz := ModuleAminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgAddPriceFeeders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
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

func (msg *MsgRemovePriceFeeders) Route() string {
	return RouterKey
}

func (msg *MsgRemovePriceFeeders) Type() string {
	return TypeMsgAddAssetInfo
}

func (msg *MsgRemovePriceFeeders) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgRemovePriceFeeders) GetSignBytes() []byte {
	bz := ModuleAminoCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgRemovePriceFeeders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid authority address (%s)", err)
	}
	return nil
}
