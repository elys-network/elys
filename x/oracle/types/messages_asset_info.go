package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateAssetInfo = "create_asset_info"
	TypeMsgUpdateAssetInfo = "update_asset_info"
	TypeMsgDeleteAssetInfo = "delete_asset_info"
)

var _ sdk.Msg = &MsgCreateAssetInfo{}

func NewMsgCreateAssetInfo(
	creator string,
	denom string,
	display string,
	bandTicker string,
	binanceTicker string,
	osmosisTicker string,
) *MsgCreateAssetInfo {
	return &MsgCreateAssetInfo{
		Creator:       creator,
		Denom:         denom,
		Display:       display,
		BandTicker:    bandTicker,
		BinanceTicker: binanceTicker,
		OsmosisTicker: osmosisTicker,
	}
}

func (msg *MsgCreateAssetInfo) Route() string {
	return RouterKey
}

func (msg *MsgCreateAssetInfo) Type() string {
	return TypeMsgCreateAssetInfo
}

func (msg *MsgCreateAssetInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateAssetInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateAssetInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateAssetInfo{}

func NewMsgUpdateAssetInfo(
	creator string,
	denom string,
	display string,
	bandTicker string,
	binanceTicker string,
	osmosisTicker string,
) *MsgUpdateAssetInfo {
	return &MsgUpdateAssetInfo{
		Creator:       creator,
		Denom:         denom,
		Display:       display,
		BandTicker:    bandTicker,
		BinanceTicker: binanceTicker,
		OsmosisTicker: osmosisTicker,
	}
}

func (msg *MsgUpdateAssetInfo) Route() string {
	return RouterKey
}

func (msg *MsgUpdateAssetInfo) Type() string {
	return TypeMsgUpdateAssetInfo
}

func (msg *MsgUpdateAssetInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateAssetInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateAssetInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeleteAssetInfo{}

func NewMsgDeleteAssetInfo(
	creator string,
	denom string,
) *MsgDeleteAssetInfo {
	return &MsgDeleteAssetInfo{
		Creator: creator,
		Denom:   denom,
	}
}
func (msg *MsgDeleteAssetInfo) Route() string {
	return RouterKey
}

func (msg *MsgDeleteAssetInfo) Type() string {
	return TypeMsgDeleteAssetInfo
}

func (msg *MsgDeleteAssetInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeleteAssetInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeleteAssetInfo) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
