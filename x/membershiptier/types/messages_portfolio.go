package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePortfolio = "create_portfolio"
	TypeMsgUpdatePortfolio = "update_portfolio"
	TypeMsgDeletePortfolio = "delete_portfolio"
)

var _ sdk.Msg = &MsgCreatePortfolio{}

func NewMsgCreatePortfolio(
	creator string,
	index string,
	assetkey string,
	coinvalue string,

) *MsgCreatePortfolio {
	return &MsgCreatePortfolio{
		Creator:   creator,
		Index:     index,
		Assetkey:  assetkey,
		Coinvalue: coinvalue,
	}
}

func (msg *MsgCreatePortfolio) Route() string {
	return RouterKey
}

func (msg *MsgCreatePortfolio) Type() string {
	return TypeMsgCreatePortfolio
}

func (msg *MsgCreatePortfolio) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePortfolio) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePortfolio) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePortfolio{}

func NewMsgUpdatePortfolio(
	creator string,
	index string,
	assetkey string,
	coinvalue string,

) *MsgUpdatePortfolio {
	return &MsgUpdatePortfolio{
		Creator:   creator,
		Index:     index,
		Assetkey:  assetkey,
		Coinvalue: coinvalue,
	}
}

func (msg *MsgUpdatePortfolio) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePortfolio) Type() string {
	return TypeMsgUpdatePortfolio
}

func (msg *MsgUpdatePortfolio) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePortfolio) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePortfolio) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePortfolio{}

func NewMsgDeletePortfolio(
	creator string,
	index string,

) *MsgDeletePortfolio {
	return &MsgDeletePortfolio{
		Creator: creator,
		Index:   index,
	}
}
func (msg *MsgDeletePortfolio) Route() string {
	return RouterKey
}

func (msg *MsgDeletePortfolio) Type() string {
	return TypeMsgDeletePortfolio
}

func (msg *MsgDeletePortfolio) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeletePortfolio) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeletePortfolio) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
