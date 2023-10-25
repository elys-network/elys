package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgOpen         = "open"
	TypeMsgClose        = "close"
	TypeMsgUpdateParams = "update_params"
	TypeMsgWhitelist    = "whitelist"
	TypeMsgUpdatePools  = "update_pools"
	TypeMsgDewhitelist  = "dewhitelist"
)

var (
	_ sdk.Msg = &MsgClose{}
	_ sdk.Msg = &MsgOpen{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgWhitelist{}
	_ sdk.Msg = &MsgUpdatePools{}
	_ sdk.Msg = &MsgDewhitelist{}
)

func NewMsgClose(creator string, id uint64) *MsgClose {
	return &MsgClose{
		Creator: creator,
		Id:      id,
	}
}

func (msg *MsgClose) Route() string {
	return RouterKey
}

func (msg *MsgClose) Type() string {
	return TypeMsgClose
}

func (msg *MsgClose) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgClose) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgOpen(creator string, collateralAsset string, collateralAmount sdk.Int, ammPoolId uint64, leverage sdk.Dec) *MsgOpen {
	return &MsgOpen{
		Creator:          creator,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		AmmPoolId:        ammPoolId,
		Leverage:         leverage,
	}
}

func (msg *MsgOpen) Route() string {
	return RouterKey
}

func (msg *MsgOpen) Type() string {
	return TypeMsgOpen
}

func (msg *MsgOpen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgOpen) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func (msg *MsgUpdateParams) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateParams) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateParams) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgUpdateParams(signer string, params *Params) *MsgUpdateParams {
	return &MsgUpdateParams{
		Authority: signer,
		Params:    params,
	}
}

func (msg *MsgUpdateParams) Route() string {
	return RouterKey
}

func (msg *MsgUpdateParams) Type() string {
	return TypeMsgUpdateParams
}

func NewMsgWhitelist(signer string, whitelistedAddress string) *MsgWhitelist {
	return &MsgWhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgWhitelist) Route() string {
	return RouterKey
}

func (msg *MsgWhitelist) Type() string {
	return TypeMsgWhitelist
}

func (msg *MsgWhitelist) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgWhitelist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgWhitelist) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgUpdatePools(signer string, pools []Pool) *MsgUpdatePools {
	return &MsgUpdatePools{
		Authority: signer,
		Pools:     pools,
	}
}

func (msg *MsgUpdatePools) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePools) Type() string {
	return TypeMsgUpdatePools
}

func (msg *MsgUpdatePools) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgUpdatePools) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePools) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

func NewMsgDewhitelist(signer string, whitelistedAddress string) *MsgDewhitelist {
	return &MsgDewhitelist{
		Authority:          signer,
		WhitelistedAddress: whitelistedAddress,
	}
}

func (msg *MsgDewhitelist) Route() string {
	return RouterKey
}

func (msg *MsgDewhitelist) Type() string {
	return TypeMsgDewhitelist
}

func (msg *MsgDewhitelist) GetSigners() []sdk.AccAddress {
	authority, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{authority}
}

func (msg *MsgDewhitelist) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDewhitelist) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Authority)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
