package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBrokerOpen = "broker_open"

func NewMsgBrokerOpen(creator string, position Position, leverage sdk.Dec, tradingAsset string, collateral sdk.Coin, takeProfitPrice sdk.Dec, owner string, stopLossPrice sdk.Dec) *MsgBrokerOpen {
	return &MsgBrokerOpen{
		Creator:         creator,
		Position:        position,
		Leverage:        leverage,
		TradingAsset:    tradingAsset,
		Collateral:      collateral,
		TakeProfitPrice: takeProfitPrice,
		Owner:           owner,
		StopLossPrice:   stopLossPrice,
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
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if msg.Position != Position_LONG && msg.Position != Position_SHORT {
		return errorsmod.Wrap(ErrInvalidPosition, msg.Position.String())
	}

	if msg.Leverage.IsNil() {
		return ErrInvalidLeverage
	}

	if !(msg.Leverage.GT(sdk.OneDec()) || msg.Leverage.IsZero()) {
		return errorsmod.Wrapf(ErrInvalidLeverage, "leverage (%s) can only be 0 (to add collateral) or > 1 to open positions", msg.Leverage.String())
	}

	if err = sdk.ValidateDenom(msg.TradingAsset); err != nil {
		return errorsmod.Wrapf(ErrInvalidTradingAsset, err.Error())
	}

	if msg.TakeProfitPrice.IsNil() {
		return errorsmod.Wrapf(ErrInvalidTakeProfitPrice, "takeProfitPrice cannot be nil")
	}

	if msg.TakeProfitPrice.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidTakeProfitPrice, "takeProfitPrice cannot be negative")
	}
	if msg.StopLossPrice.IsNil() {
		return errorsmod.Wrapf(ErrInvalidPrice, "stopLossPrice cannot be nil")
	}

	if msg.StopLossPrice.IsNegative() {
		return errorsmod.Wrapf(ErrInvalidPrice, "stopLossPrice cannot be negative")
	}

	if msg.Position == Position_LONG && !msg.StopLossPrice.IsZero() && msg.TakeProfitPrice.LTE(msg.StopLossPrice) {
		return fmt.Errorf("TakeProfitPrice cannot be <= StopLossPrice for LONG")
	}
	if msg.Position == Position_SHORT && !msg.StopLossPrice.IsZero() && msg.TakeProfitPrice.GTE(msg.StopLossPrice) {
		return fmt.Errorf("TakeProfitPrice cannot be >= StopLossPrice for SHORT")
	}

	return nil
}
