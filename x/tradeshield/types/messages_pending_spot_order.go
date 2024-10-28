package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateSpotOrder  = "create_spot_order"
	TypeMsgUpdateSpotOrder  = "update_spot_order"
	TypeMsgCancelSpotOrders = "cancel_spot_orders"
)

var _ sdk.Msg = &MsgCreateSpotOrder{}

func NewMsgCreateSpotOrder(ownerAddress string, orderType SpotOrderType,
	orderPrice OrderPrice, orderAmount sdk.Coin,
	orderTargetDenom string) *MsgCreateSpotOrder {
	return &MsgCreateSpotOrder{
		OrderType:        orderType,
		OrderPrice:       &orderPrice,
		OrderAmount:      &orderAmount,
		OwnerAddress:     ownerAddress,
		OrderTargetDenom: orderTargetDenom,
	}
}

func (msg *MsgCreateSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgCreateSpotOrder) Type() string {
	return TypeMsgCreateSpotOrder
}

func (msg *MsgCreateSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreateSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreateSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateSpotOrder{}

func NewMsgUpdateSpotOrder(creator string, id uint64, orderPrice *OrderPrice) *MsgUpdateSpotOrder {
	return &MsgUpdateSpotOrder{
		OrderId:      id,
		OwnerAddress: creator,
		OrderPrice:   orderPrice,
	}
}

func (msg *MsgUpdateSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgUpdateSpotOrder) Type() string {
	return TypeMsgUpdateSpotOrder
}

func (msg *MsgUpdateSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdateSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdateSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgCancelSpotOrders{}

func NewMsgCancelSpotOrders(creator string, id []uint64) *MsgCancelSpotOrders {
	return &MsgCancelSpotOrders{
		SpotOrderIds: id,
		Creator:      creator,
	}
}
func (msg *MsgCancelSpotOrders) Route() string {
	return RouterKey
}

func (msg *MsgCancelSpotOrders) Type() string {
	return TypeMsgCancelSpotOrders
}

func (msg *MsgCancelSpotOrders) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelSpotOrders) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelSpotOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}
