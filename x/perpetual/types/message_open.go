package types

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgOpen{}

func NewMsgOpen(creator string, position Position, leverage sdkmath.LegacyDec, poolId uint64, collateral sdk.Coin, takeProfitPrice sdkmath.LegacyDec, stopLossPrice sdkmath.LegacyDec) *MsgOpen {
	return &MsgOpen{
		Creator:         creator,
		Position:        position,
		Leverage:        leverage,
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

	if err = msg.Collateral.Validate(); err != nil {
		return err
	}
	if err = CheckLegacyDecNilAndNegative(msg.TakeProfitPrice, "TakeProfitPrice"); err != nil {
		return err
	}
	if err = CheckLegacyDecNilAndNegative(msg.StopLossPrice, "StopLossPrice"); err != nil {
		return err
	}
	if msg.Position == Position_LONG && !msg.StopLossPrice.IsZero() && !msg.TakeProfitPrice.IsZero() && msg.TakeProfitPrice.LTE(msg.StopLossPrice) {
		return errors.New("TakeProfitPrice cannot be <= StopLossPrice for LONG")
	}
	if msg.Position == Position_SHORT && !msg.StopLossPrice.IsZero() && !msg.TakeProfitPrice.IsZero() && msg.TakeProfitPrice.GTE(msg.StopLossPrice) {
		return errors.New("TakeProfitPrice cannot be >= StopLossPrice for SHORT")
	}
	return nil
}

func (msg *MsgOpen) GetCustodyAndLiabilitiesAsset(tradingAsset, baseCurrency string) (string, string, error) {
	// Define the assets
	var liabilitiesAsset, custodyAsset string
	switch msg.Position {
	case Position_LONG:
		liabilitiesAsset = baseCurrency
		custodyAsset = tradingAsset
	case Position_SHORT:
		liabilitiesAsset = tradingAsset
		custodyAsset = baseCurrency
	default:
		return "", "", errorsmod.Wrap(ErrInvalidPosition, msg.Position.String())
	}
	return custodyAsset, liabilitiesAsset, nil
}

func (msg *MsgOpen) ValidatePosition(tradingAsset, baseCurrency string) error {
	// Determine the type of position (long or short) and validate assets accordingly.
	switch msg.Position {
	case Position_LONG:
		if msg.Collateral.Denom != tradingAsset && msg.Collateral.Denom != baseCurrency {
			return errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid operation: collateral asset has to be either trading asset or base currency for long")
		}
	case Position_SHORT:
		// The collateral for a short must be the base currency.
		if msg.Collateral.Denom != baseCurrency {
			return errorsmod.Wrap(ErrInvalidCollateralAsset, "invalid collateral: collateral asset for short position must be the base currency")
		}
	default:
		return errorsmod.Wrap(ErrInvalidPosition, msg.Position.String())
	}
	return nil
}

func (msg *MsgOpen) ValidateTakeProfitAndStopLossPrice(params Params, tradingAssetPrice sdkmath.LegacyDec) error {
	// if msg.TakeProfitPrice is not positive, then it has to be 0 because of validate basic
	if msg.TakeProfitPrice.IsPositive() {
		ratio := msg.TakeProfitPrice.Quo(tradingAssetPrice)
		if msg.Position == Position_LONG {
			if ratio.LT(params.MinimumLongTakeProfitPriceRatio) || ratio.GT(params.MaximumLongTakeProfitPriceRatio) {
				return fmt.Errorf("take profit price should be between %s and %s times of current market price for long (current ratio: %s)", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String(), ratio.String())
			}
			if !msg.StopLossPrice.IsZero() && msg.StopLossPrice.GTE(tradingAssetPrice) {
				return fmt.Errorf("stop loss price cannot be greater than equal to tradingAssetPrice for long (Stop loss: %s, asset price: %s)", msg.StopLossPrice.String(), tradingAssetPrice.String())
			}
			// no need to override msg.TakeProfitPrice as the above ratio checks it
		}
		if msg.Position == Position_SHORT {
			if ratio.GT(params.MaximumShortTakeProfitPriceRatio) {
				return fmt.Errorf("take profit price should be less than %s times of current market price for short (current ratio: %s)", params.MaximumShortTakeProfitPriceRatio.String(), ratio.String())
			}
			if !msg.StopLossPrice.IsZero() && msg.StopLossPrice.LTE(tradingAssetPrice) {
				return fmt.Errorf("stop loss price cannot be less than equal to tradingAssetPrice for short (Stop loss: %s, asset price: %s)", msg.StopLossPrice.String(), tradingAssetPrice.String())
			}
		}
	}
	return nil
}
