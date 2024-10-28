package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePendingSpotOrder = "create_pending_spot_order"
	TypeMsgUpdatePendingSpotOrder = "update_pending_spot_order"
	TypeMsgCancelSpotOrders       = "cancel_pending_spot_order"
)

var _ sdk.Msg = &MsgCreatePendingSpotOrder{}

func NewMsgCreatePendingSpotOrder(ownerAddress string, orderType SpotOrderType,
	orderPrice OrderPrice, orderAmount sdk.Coin,
	orderTargetDenom string) *MsgCreatePendingSpotOrder {
	return &MsgCreatePendingSpotOrder{
		OrderType:        orderType,
		OrderPrice:       &orderPrice,
		OrderAmount:      &orderAmount,
		OwnerAddress:     ownerAddress,
		OrderTargetDenom: orderTargetDenom,
	}
}

func (msg *MsgCreatePendingSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgCreatePendingSpotOrder) Type() string {
	return TypeMsgCreatePendingSpotOrder
}

func (msg *MsgCreatePendingSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePendingSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePendingSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePendingSpotOrder{}

func NewMsgUpdatePendingSpotOrder(creator string, id uint64, orderPrice *OrderPrice) *MsgUpdatePendingSpotOrder {
	return &MsgUpdatePendingSpotOrder{
		OrderId:      id,
		OwnerAddress: creator,
		OrderPrice:   orderPrice,
	}
}

func (msg *MsgUpdatePendingSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePendingSpotOrder) Type() string {
	return TypeMsgUpdatePendingSpotOrder
}

func (msg *MsgUpdatePendingSpotOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePendingSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePendingSpotOrder) ValidateBasic() error {
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
