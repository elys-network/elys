package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgCreateAssetInfo = "create_asset_info"

var _ sdk.Msg = &MsgCreateAssetInfo{}

func NewMsgCreateAssetInfo(creator string, denom string, display string, bandTicker string, elysTicker string, decimal uint64) *MsgCreateAssetInfo {
	return &MsgCreateAssetInfo{
		Creator:    creator,
		Denom:      denom,
		Display:    display,
		BandTicker: bandTicker,
		ElysTicker: elysTicker,
		Decimal:    decimal,
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
