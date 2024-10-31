package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreateSpotOrder  = "create_spot_order"
	TypeMsgUpdateSpotOrder  = "update_spot_order"
	TypeMsgCancelSpotOrder  = "cancel_spot_order"
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

	// Validate order price
	if msg.OrderPrice == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be nil")
	}

	// Validate order amount
	if !msg.OrderAmount.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid order amount")
	}

	// Validate order target denom
	if msg.OrderTargetDenom == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order target denom cannot be empty")
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

	// Validate order price
	if msg.OrderPrice == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be nil")
	}

	// Validate order price
	if msg.OrderPrice.Rate.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be negative")
	}

	err = sdk.ValidateDenom(msg.OrderPrice.BaseDenom)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid base asset denom (%s)", err)
	}

	err = sdk.ValidateDenom(msg.OrderPrice.QuoteDenom)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "invalid quote asset denom (%s)", err)
	}

	// Validate order ID
	if msg.OrderId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be 0")
	}
	return nil
}

var _ sdk.Msg = &MsgCancelSpotOrder{}

func NewMsgCancelSpotOrder(ownerAddress string, orderId uint64) *MsgCancelSpotOrder {
	return &MsgCancelSpotOrder{
		OwnerAddress: ownerAddress,
		OrderId:      orderId,
	}
}

func (msg *MsgCancelSpotOrder) Route() string {
	return RouterKey
}

func (msg *MsgCancelSpotOrder) Type() string {
	return TypeMsgCancelSpotOrder
}

func (msg *MsgCancelSpotOrder) GetSigners() []sdk.AccAddress {
	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{ownerAddress}
}

func (msg *MsgCancelSpotOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ownerAddress address (%s)", err)
	}

	// Validate order ID
	if msg.OrderId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be 0")
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

	// Validate SpotOrderIds
	if len(msg.SpotOrderIds) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "spot order IDs cannot be empty")
	}
	for _, id := range msg.SpotOrderIds {
		if id == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
		}
	}
	return nil
}
