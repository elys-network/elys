package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgOpen         = "open"
	TypeMsgClose        = "close"
	TypeMsgBrokerOpen   = "broker_open"
	TypeMsgBrokerClose  = "broker_close"
	TypeMsgUpdateParams = "update_params"
	TypeMsgWhitelist    = "whitelist"
	TypeMsgUpdatePools  = "update_pools"
	TypeMsgDewhitelist  = "dewhitelist"
)

var (
	_ sdk.Msg = &MsgClose{}
	_ sdk.Msg = &MsgOpen{}
	_ sdk.Msg = &MsgBrokerClose{}
	_ sdk.Msg = &MsgBrokerOpen{}
	_ sdk.Msg = &MsgUpdateParams{}
	_ sdk.Msg = &MsgWhitelist{}
	_ sdk.Msg = &MsgUpdatePools{}
	_ sdk.Msg = &MsgDewhitelist{}
)

func NewMsgClose(creator string, id uint64, amount sdk.Coin) *MsgClose {
	return &MsgClose{
		Creator: creator,
		Id:      id,
		Amount:  amount,
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

func NewMsgOpen(creator string, collateralAsset string, collateralAmount sdk.Int, borrowAsset string, position Position, leverage sdk.Dec, takeProfitPrice sdk.Dec) *MsgOpen {
	return &MsgOpen{
		Creator:          creator,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		BorrowAsset:      borrowAsset,
		Position:         position,
		Leverage:         leverage,
		TakeProfitPrice:  takeProfitPrice,
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

func NewMsgBrokerClose(creator string, id uint64, owner string, amount sdk.Coin) *MsgBrokerClose {
	return &MsgBrokerClose{
		Creator: creator,
		Id:      id,
		Owner:   owner,
		Amount:  amount,
	}
}

func (msg *MsgBrokerClose) Route() string {
	return RouterKey
}

func (msg *MsgBrokerClose) Type() string {
	return TypeMsgBrokerClose
}

func (msg *MsgBrokerClose) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBrokerClose) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBrokerClose) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}

func NewMsgBrokerOpen(creator string, collateralAsset string, collateralAmount sdk.Int, borrowAsset string, position Position, leverage sdk.Dec, takeProfitPrice sdk.Dec, owner string) *MsgBrokerOpen {
	return &MsgBrokerOpen{
		Creator:          creator,
		CollateralAsset:  collateralAsset,
		CollateralAmount: collateralAmount,
		BorrowAsset:      borrowAsset,
		Position:         position,
		Leverage:         leverage,
		TakeProfitPrice:  takeProfitPrice,
		Owner:            owner,
	}
}

func (msg *MsgBrokerOpen) Route() string {
	return RouterKey
}

func (msg *MsgBrokerOpen) Type() string {
	return TypeMsgBrokerOpen
}

func (msg *MsgBrokerOpen) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgBrokerOpen) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgBrokerOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
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

func NewMsgUpdatePools(signer string, pools []string, closedPools []string) *MsgUpdatePools {
	return &MsgUpdatePools{
		Authority:   signer,
		Pools:       pools,
		ClosedPools: closedPools,
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
