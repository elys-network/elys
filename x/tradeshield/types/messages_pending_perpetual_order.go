package types

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatePerpetualOpenOrder{}

func NewMsgCreatePerpetualOpenOrder(
	ownerAddress string,
	triggerPrice TriggerPrice,
	collateral sdk.Coin,
	tradingAsset string,
	position PerpetualPosition,
	leverage math.LegacyDec,
	takeProfitPrice math.LegacyDec,
	stopLossPrice math.LegacyDec,
	poolId uint64,
) *MsgCreatePerpetualOpenOrder {
	return &MsgCreatePerpetualOpenOrder{
		TriggerPrice:    triggerPrice,
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

func (msg *MsgCreatePerpetualOpenOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err = CheckLegacyDecNilAndNegative(msg.TriggerPrice.Rate, "TriggerPrice Rate"); err != nil {
		return err
	}

	if err = sdk.ValidateDenom(msg.TriggerPrice.TradingAssetDenom); err != nil {
		return err
	}

	// Validate collateral
	if !msg.Collateral.IsValid() {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid collateral")
	}

	if err = sdk.ValidateDenom(msg.TradingAsset); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid trading asset denom (%s)", err)
	}

	if msg.Position != PerpetualPosition_LONG && msg.Position != PerpetualPosition_SHORT {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "invalid position")
	}

	if err = CheckLegacyDecNilAndNegative(msg.Leverage, "Leverage"); err != nil {
		return err
	}

	if err = CheckLegacyDecNilAndNegative(msg.TakeProfitPrice, "TakeProfitPrice"); err != nil {
		return err
	}

	if err = CheckLegacyDecNilAndNegative(msg.StopLossPrice, "StopLossPrice"); err != nil {
		return err
	}

	// Validate pool ID
	if msg.PoolId == 0 {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "pool ID cannot be zero")
	}

	if msg.Position == PerpetualPosition_LONG && !msg.StopLossPrice.IsZero() && msg.TakeProfitPrice.LTE(msg.StopLossPrice) {
		return fmt.Errorf("TakeProfitPrice cannot be <= StopLossPrice for LONG")
	}
	if msg.Position == PerpetualPosition_SHORT && !msg.StopLossPrice.IsZero() && msg.TakeProfitPrice.GTE(msg.StopLossPrice) {
		return fmt.Errorf("TakeProfitPrice cannot be >= StopLossPrice for SHORT")
	}
	return nil
}

var _ sdk.Msg = &MsgCreatePerpetualOpenOrder{}

func NewMsgCreatePerpetualCloseOrder(
	ownerAddress string,
	triggerPrice TriggerPrice,
	positionId uint64,
) *MsgCreatePerpetualCloseOrder {
	return &MsgCreatePerpetualCloseOrder{
		TriggerPrice: triggerPrice,
		OwnerAddress: ownerAddress,
		PositionId:   positionId,
	}
}

func (msg *MsgCreatePerpetualCloseOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if err = CheckLegacyDecNilAndNegative(msg.TriggerPrice.Rate, "TriggerPrice Rate"); err != nil {
		return err
	}

	err = sdk.ValidateDenom(msg.TriggerPrice.TradingAssetDenom)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid trading asset denom (%s)", err)
	}

	// Validate PositionId
	if msg.PositionId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "position ID cannot be zero")
	}
	return nil
}

var _ sdk.Msg = &MsgUpdatePerpetualOrder{}

func NewMsgUpdatePerpetualOrder(creator string, id uint64, triggerPrice TriggerPrice) *MsgUpdatePerpetualOrder {
	return &MsgUpdatePerpetualOrder{
		OrderId:      id,
		OwnerAddress: creator,
		TriggerPrice: triggerPrice,
	}
}

func (msg *MsgUpdatePerpetualOrder) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if err = CheckLegacyDecNilAndNegative(msg.TriggerPrice.Rate, "TriggerPrice Rate"); err != nil {
		return err
	}

	err = sdk.ValidateDenom(msg.TriggerPrice.TradingAssetDenom)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid trading asset denom (%s)", err)
	}

	// Validate Order ID
	if msg.OrderId == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "Order ID cannot be zero")
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

func (msg *MsgCancelPerpetualOrder) ValidateBasic() error {
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

var _ sdk.Msg = &MsgCancelPerpetualOrders{}

func NewMsgCancelPerpetualOrders(creator string, ids []uint64) *MsgCancelPerpetualOrders {
	return &MsgCancelPerpetualOrders{
		OrderIds:     ids,
		OwnerAddress: creator,
	}
}

func (msg *MsgCancelPerpetualOrders) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.OwnerAddress)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	// Validate SpotOrderIds
	if len(msg.OrderIds) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "spot order IDs cannot be empty")
	}
	for _, id := range msg.OrderIds {
		if id == 0 {
			return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "spot order ID cannot be zero")
		}
	}

	return nil
}
