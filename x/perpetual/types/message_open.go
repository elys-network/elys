package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgOpen{}

func NewMsgOpen(creator string, position Position, leverage sdkmath.LegacyDec, poolId uint64, tradingAsset string, collateral sdk.Coin, takeProfitPrice sdkmath.LegacyDec, stopLossPrice sdkmath.LegacyDec) *MsgOpen {
	return &MsgOpen{
		Creator:         creator,
		Position:        position,
		Leverage:        leverage,
		TradingAsset:    tradingAsset,
		Collateral:      collateral,
		TakeProfitPrice: takeProfitPrice,
		StopLossPrice:   stopLossPrice,
		PoolId:          poolId,
	}
}

func (msg *MsgOpen) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid creator address (%s)", err)
	}

	if msg.Position.String() != "LONG" && msg.Position.String() != "SHORT" {
		return errorsmod.Wrap(ErrInvalidPosition, msg.Position.String())
	}

	if msg.Leverage.IsNil() {
		return ErrInvalidLeverage
	}

	if msg.Leverage.LT(sdk.OneDec()) {
		return errorsmod.Wrapf(ErrInvalidLeverage, "leverage (%s) cannot be < 1", msg.Leverage.String())
	}

	if len(msg.TradingAsset) == 0 {
		return ErrTradingAssetIsEmpty
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
	return nil
}
