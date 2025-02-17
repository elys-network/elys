package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreateSpotOrder{}

func NewMsgCreateSpotOrder(ownerAddress string, orderType SpotOrderType,
	orderPrice math.LegacyDec, orderAmount sdk.Coin,
	orderTargetDenom string) *MsgCreateSpotOrder {
	return &MsgCreateSpotOrder{
		OrderType:        orderType,
		OrderPrice:       orderPrice,
		OrderAmount:      orderAmount,
		OwnerAddress:     ownerAddress,
		OrderTargetDenom: orderTargetDenom,
	}
}

func (msg *MsgCreateSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err = CheckLegacyDecNilAndNegative(msg.OrderPrice, "OrderPrice Rate"); err != nil {
		return err
	}

	if err = sdk.ValidateDenom(msg.OrderAmount.Denom); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid order amount denom (%s)", err)
	}

	// check that order amount denom is not the same as the order target denom
	if msg.OrderAmount.Denom == msg.OrderTargetDenom {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "order amount denom cannot be the same as the order target denom")
	}

	// Validate order amount
	if !msg.OrderAmount.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid order amount")
	}

	if err = sdk.ValidateDenom(msg.OrderTargetDenom); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid order target asset denom (%s)", err)
	}

	if msg.OrderType != SpotOrderType_STOPLOSS && msg.OrderType != SpotOrderType_LIMITSELL && msg.OrderType != SpotOrderType_LIMITBUY && msg.OrderType != SpotOrderType_MARKETBUY {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid order type")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdateSpotOrder{}

func NewMsgUpdateSpotOrder(creator string, id uint64, orderPrice math.LegacyDec) *MsgUpdateSpotOrder {
	return &MsgUpdateSpotOrder{
		OrderId:      id,
		OwnerAddress: creator,
		OrderPrice:   orderPrice,
	}
}

func (msg *MsgUpdateSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err = CheckLegacyDecNilAndNegative(msg.OrderPrice, "OrderPrice Rate"); err != nil {
		return err
	}

	// Validate order ID
	if msg.OrderId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be 0")
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

func (msg *MsgCancelSpotOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid ownerAddress address (%s)", err)
	}

	// Validate order ID
	if msg.OrderId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "order price cannot be 0")
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

func (msg *MsgCancelSpotOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// Validate SpotOrderIds
	if len(msg.SpotOrderIds) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "spot order IDs cannot be empty")
	}
	for _, id := range msg.SpotOrderIds {
		if id == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
		}
	}
	return nil
}
