package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePrice = "create_price"
	TypeMsgUpdatePrice = "update_price"
	TypeMsgDeletePrice = "delete_price"
)

var _ sdk.Msg = &MsgCreatePrice{}

func NewMsgCreatePrice(
	creator string,
	asset string,
	price sdk.Dec,
	source string,
) *MsgCreatePrice {
	return &MsgCreatePrice{
		Provider: creator,
		Asset:    asset,
		Price:    price,
		Source:   source,
	}
}

func (msg *MsgCreatePrice) Route() string {
	return RouterKey
}

func (msg *MsgCreatePrice) Type() string {
	return TypeMsgCreatePrice
}

func (msg *MsgCreatePrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePrice{}

func NewMsgUpdatePrice(
	creator string,
	asset string,
	price sdk.Dec,
	source string,
) *MsgUpdatePrice {
	return &MsgUpdatePrice{
		Provider: creator,
		Asset:    asset,
		Price:    price,
		Source:   source,
	}
}

func (msg *MsgUpdatePrice) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePrice) Type() string {
	return TypeMsgUpdatePrice
}

func (msg *MsgUpdatePrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Provider)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePrice{}

func NewMsgDeletePrice(
	creator string,
	asset string,
) *MsgDeletePrice {
	return &MsgDeletePrice{
		Creator: creator,
		Asset:   asset,
	}
}
func (msg *MsgDeletePrice) Route() string {
	return RouterKey
}

func (msg *MsgDeletePrice) Type() string {
	return TypeMsgDeletePrice
}

func (msg *MsgDeletePrice) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeletePrice) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeletePrice) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
