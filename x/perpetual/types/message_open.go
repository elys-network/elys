package types

import (
	"errors"

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

	if msg.Position != Position_LONG && msg.Position != Position_SHORT {
		return errorsmod.Wrap(ErrInvalidPosition, msg.Position.String())
	}

	if msg.Leverage.IsNil() {
		return ErrInvalidLeverage
	}

	if msg.PoolId == 0 {
		return errors.New("pool id cannot be 0")
	}

	if !(msg.Leverage.GT(sdkmath.LegacyOneDec()) || msg.Leverage.IsZero()) {
		return errorsmod.Wrapf(ErrInvalidLeverage, "leverage (%s) can only be 0 (to add collateral) or > 1 to open positions", msg.Leverage.String())
	}

	if err = sdk.ValidateDenom(msg.TradingAsset); err != nil {
		return errorsmod.Wrapf(ErrInvalidTradingAsset, err.Error())
	}

	if err = msg.Collateral.Validate(); err != nil {
		return err
	}
	if err = CheckLegacyDecNilAndNegative(msg.TakeProfitPrice, "TakeProfitPrice"); err != nil {
		return err
	}
	if err = CheckLegacyDecNilAndNegative(msg.StopLossPrice, "StopLossPrice"); err != nil {
		return err
	}
	if msg.Position == Position_LONG && !msg.StopLossPrice.IsZero() && msg.TakeProfitPrice.LTE(msg.StopLossPrice) {
		return errors.New("TakeProfitPrice cannot be <= StopLossPrice for LONG")
	}
	if msg.Position == Position_SHORT && !msg.StopLossPrice.IsZero() && msg.TakeProfitPrice.GTE(msg.StopLossPrice) {
		return errors.New("TakeProfitPrice cannot be >= StopLossPrice for SHORT")
	}
	return nil
}
