package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePendingPerpetualOrder = "create_pending_perpetual_order"
	TypeMsgUpdatePendingPerpetualOrder = "update_pending_perpetual_order"
	TypeMsgDeletePendingPerpetualOrder = "delete_pending_perpetual_order"
)

var _ sdk.Msg = &MsgCreatePendingPerpetualOrder{}

func NewMsgCreatePendingPerpetualOrder(creator string, order string) *MsgCreatePendingPerpetualOrder {
	return &MsgCreatePendingPerpetualOrder{
		Creator: creator,
		Order:   order,
	}
}

func (msg *MsgCreatePendingPerpetualOrder) Route() string {
	return RouterKey
}

func (msg *MsgCreatePendingPerpetualOrder) Type() string {
	return TypeMsgCreatePendingPerpetualOrder
}

func (msg *MsgCreatePendingPerpetualOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePendingPerpetualOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePendingPerpetualOrder{}

func NewMsgUpdatePendingPerpetualOrder(creator string, id uint64, order string) *MsgUpdatePendingPerpetualOrder {
	return &MsgUpdatePendingPerpetualOrder{
		Id:      id,
		Creator: creator,
		Order:   order,
	}
}

func (msg *MsgUpdatePendingPerpetualOrder) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePendingPerpetualOrder) Type() string {
	return TypeMsgUpdatePendingPerpetualOrder
}

func (msg *MsgUpdatePendingPerpetualOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePendingPerpetualOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgDeletePendingPerpetualOrder{}

func NewMsgDeletePendingPerpetualOrder(creator string, id uint64) *MsgDeletePendingPerpetualOrder {
	return &MsgDeletePendingPerpetualOrder{
		Id:      id,
		Creator: creator,
	}
}
func (msg *MsgDeletePendingPerpetualOrder) Route() string {
	return RouterKey
}

func (msg *MsgDeletePendingPerpetualOrder) Type() string {
	return TypeMsgDeletePendingPerpetualOrder
}

func (msg *MsgDeletePendingPerpetualOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgDeletePendingPerpetualOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgDeletePendingPerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}
