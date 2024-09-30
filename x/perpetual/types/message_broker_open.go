package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const TypeMsgBrokerOpen = "broker_open"

func NewMsgBrokerOpen(creator string, position Position, leverage sdkmath.LegacyDec, tradingAsset string, collateral sdk.Coin, takeProfitPrice sdkmath.LegacyDec, owner string, stopLossPrice sdkmath.LegacyDec) *MsgBrokerOpen {
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

func (msg *MsgBrokerOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}
	_, err = sdk.AccAddressFromBech32(msg.Owner)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid owner address (%s)", err)
	}

	if msg.Position.String() != "LONG" && msg.Position.String() != "SHORT" {
		return errorsmod.Wrap(ErrInvalidPosition, msg.Position.String())
	}

	if msg.Leverage.IsNil() {
		return ErrInvalidLeverage
	}

	if msg.Leverage.IsNegative() {
		return ErrInvalidLeverage
	}

	if len(msg.TradingAsset) == 0 {
		return ErrTradingAssetIsEmpty
	}

	if msg.TakeProfitPrice.IsNil() {
		return ErrInvalidTakeProfitPriceIsNegative
	}

	if msg.TakeProfitPrice.IsNegative() {
		return ErrInvalidTakeProfitPriceIsNegative
	}

	return nil
}
