package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgSetAssetInfo    = "update_asset_info"
	TypeMsgDeleteAssetInfo = "delete_asset_info"
)

var _ sdk.Msg = &MsgSetAssetInfo{}

func NewMsgSetAssetInfo(
	creator string,
	denom string,
	display string,
	bandTicker string,
	elysTicker string,
) *MsgSetAssetInfo {
	return &MsgSetAssetInfo{
		Creator:    creator,
		Denom:      denom,
		Display:    display,
		BandTicker: bandTicker,
		ElysTicker: elysTicker,
	}
}

func (msg *MsgSetAssetInfo) Route() string {
	return RouterKey
}

func (msg *MsgSetAssetInfo) Type() string {
	return TypeMsgSetAssetInfo
}

func (msg *MsgSetAssetInfo) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgSetAssetInfo) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgSetAssetInfo) ValidateBasic() error {
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
