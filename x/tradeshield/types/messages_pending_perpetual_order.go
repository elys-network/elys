package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	TypeMsgCreatePerpetualOpenOrder  = "create_perpetual_open_order"
	TypeMsgCreatePerpetualCloseOrder = "create_perpetual_close_order"
	TypeMsgUpdatePerpetualOrder      = "update_perpetual_order"
	TypeMsgCancelPerpetualOrder      = "cancel_perpetual_order"
	TypeMsgCancelPerpetualOrders     = "cancel_perpetual_orders"
)

var _ sdk.Msg = &MsgCreatePerpetualOpenOrder{}

func NewMsgCreatePerpetualOpenOrder(
	ownerAddress string,
	orderType PerpetualOrderType,
	triggerPrice TriggerPrice,
	collateral sdk.Coin,
	tradingAsset string,
	position PerpetualPosition,
	leverage sdk.Dec,
	takeProfitPrice sdk.Dec,
	stopLossPrice sdk.Dec,
	poolId uint64,
) *MsgCreatePerpetualOpenOrder {
	return &MsgCreatePerpetualOpenOrder{
		TriggerPrice:    &triggerPrice,
		Collateral:      collateral,
		OwnerAddress:    ownerAddress,
		TradingAsset:    tradingAsset,
		Position:        position,
		Leverage:        leverage,
		TakeProfitPrice: takeProfitPrice,
		StopLossPrice:   stopLossPrice,
		PoolId:          poolId,
	}
}

func (msg *MsgCreatePerpetualOpenOrder) Route() string {
	return RouterKey
}

func (msg *MsgCreatePerpetualOpenOrder) Type() string {
	return TypeMsgCreatePerpetualOpenOrder
}

func (msg *MsgCreatePerpetualOpenOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePerpetualOpenOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePerpetualOpenOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// Validate trigger price
	if msg.TriggerPrice == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trigger price cannot be nil")
	}

	// Validate collateral
	if !msg.Collateral.IsValid() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "invalid collateral")
	}

	// Validate trading asset
	if msg.TradingAsset == "" {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trading asset cannot be empty")
	}

	// Validate leverage
	if msg.Leverage.IsNil() || msg.Leverage.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "leverage cannot be nil or negative")
	}

	// Validate take profit price
	if msg.TakeProfitPrice.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "take profit price cannot be negative")
	}

	// Validate stop loss price
	if msg.StopLossPrice.IsNegative() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "stop loss price cannot be negative")
	}

	// Validate pool ID
	if msg.PoolId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "pool ID cannot be zero")
	}

	return nil
}

var _ sdk.Msg = &MsgCreatePerpetualOpenOrder{}

func NewMsgCreatePerpetualCloseOrder(
	ownerAddress string,
	orderType PerpetualOrderType,
	triggerPrice TriggerPrice,
	positionId uint64,
) *MsgCreatePerpetualCloseOrder {
	return &MsgCreatePerpetualCloseOrder{
		TriggerPrice: &triggerPrice,
		OwnerAddress: ownerAddress,
		PositionId:   positionId,
	}
}

func (msg *MsgCreatePerpetualCloseOrder) Route() string {
	return RouterKey
}

func (msg *MsgCreatePerpetualCloseOrder) Type() string {
	return TypeMsgCreatePerpetualOpenOrder
}

func (msg *MsgCreatePerpetualCloseOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCreatePerpetualCloseOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCreatePerpetualCloseOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePerpetualOrder{}

func NewMsgUpdatePerpetualOrder(creator string, id uint64, triggerPrice *TriggerPrice) *MsgUpdatePerpetualOrder {
	return &MsgUpdatePerpetualOrder{
		OrderId:      id,
		OwnerAddress: creator,
		TriggerPrice: triggerPrice,
	}
}

func (msg *MsgUpdatePerpetualOrder) Route() string {
	return RouterKey
}

func (msg *MsgUpdatePerpetualOrder) Type() string {
	return TypeMsgUpdatePerpetualOrder
}

func (msg *MsgUpdatePerpetualOrder) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgUpdatePerpetualOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgUpdatePerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	// Validate trigger price
	if msg.TriggerPrice == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "trigger price cannot be nil")
	}

	// Validate Order ID
	if msg.OrderId == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Order ID cannot be zero")
	}
	return nil
}

var _ sdk.Msg = &MsgCancelPerpetualOrder{}

func NewMsgCancelPerpetualOrder(ownerAddress string, orderId uint64) *MsgCancelPerpetualOrder {
	return &MsgCancelPerpetualOrder{
		OwnerAddress: ownerAddress,
		OrderId:      orderId,
	}
}

func (msg *MsgCancelPerpetualOrder) Route() string {
	return RouterKey
}

func (msg *MsgCancelPerpetualOrder) Type() string {
	return TypeMsgCancelPerpetualOrder
}

func (msg *MsgCancelPerpetualOrder) GetSigners() []sdk.AccAddress {
	ownerAddress, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{ownerAddress}
}

func (msg *MsgCancelPerpetualOrder) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelPerpetualOrder) ValidateBasic() error {
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

var _ sdk.Msg = &MsgCancelPerpetualOrders{}

func NewMsgCancelPerpetualOrders(creator string, ids []uint64) *MsgCancelPerpetualOrders {
	return &MsgCancelPerpetualOrders{
		OrderIds:     ids,
		OwnerAddress: creator,
	}
}
func (msg *MsgCancelPerpetualOrders) Route() string {
	return RouterKey
}

func (msg *MsgCancelPerpetualOrders) Type() string {
	return TypeMsgCancelPerpetualOrders
}

func (msg *MsgCancelPerpetualOrders) GetSigners() []sdk.AccAddress {
	creator, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{creator}
}

func (msg *MsgCancelPerpetualOrders) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg *MsgCancelPerpetualOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// Validate SpotOrderIds
	if len(msg.OrderIds) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "spot order IDs cannot be empty")
	}
	for _, id := range msg.OrderIds {
		if id == 0 {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
		}
	}

	return nil
}
