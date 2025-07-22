package types

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var _ sdk.Msg = &MsgCreatPerpetualMarket{}
var _ sdk.Msg = &MsgWithdraw{}

func (msg MsgCreatPerpetualMarket) ValidateBasic() (err error) {
	err = sdk.ValidateDenom(msg.BaseDenom)
	if err != nil {
		return err
	}
	err = sdk.ValidateDenom(msg.QuoteDenom)
	if err != nil {
		return err
	}
	if msg.TwapPricesWindow <= 10 {
		return fmt.Errorf("max twap prices time must be greater than 10")
	}
	if msg.MaxAbsFundingRate.IsNil() || msg.MaxAbsFundingRate.IsNegative() {
		return fmt.Errorf("max abs funding rate cannot be negative or nil")
	}
	if msg.MaxAbsFundingRateChange.IsNil() || msg.MaxAbsFundingRateChange.IsNegative() {
		return fmt.Errorf("max abs funding rate cannot be negative or nil")
	}
	if msg.InitialMarginRatio.IsNil() || msg.InitialMarginRatio.IsNegative() || msg.InitialMarginRatio.IsZero() {
		return fmt.Errorf("initial margin ratio cannot be negative or zero")
	}
	if msg.MaintenanceMarginRatio.IsNil() || msg.MaintenanceMarginRatio.IsNegative() || msg.MaintenanceMarginRatio.IsZero() {
		return fmt.Errorf("maintenance margin ratio cannot be negative or zero")
	}
	if msg.MaintenanceMarginRatio.GTE(msg.InitialMarginRatio) {
		return fmt.Errorf("maintenance margin ratio cannot be greater than or equal to initial margin ratio")
	}
	imrMMRRatio := msg.InitialMarginRatio.Quo(msg.MaintenanceMarginRatio)
	if imrMMRRatio.LT(MinimumIMRByMMR) || imrMMRRatio.GT(MaximumIMRByMMR) {
		return fmt.Errorf("imr/mmr must be between %s and %s", MinimumIMRByMMR.String(), MaximumIMRByMMR.String())
	}
	return nil
}

func (msg MsgDeposit) ValidateBasic() error {
	// Validate sender address
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address: %s", err)
	}

	// Validate coin
	if err := msg.Coin.Validate(); err != nil {
		return errorsmod.Wrapf(ErrInvalidCoin, "invalid coin: %s", err)
	}

	// Check for positive amount
	if msg.Coin.Amount.IsNil() || !msg.Coin.Amount.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidCoin, "coin amount must be positive")
	}

	return nil
}

func (msg MsgPlaceLimitOrder) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Validate market ID (must be non-zero)
	if msg.MarketId == 0 {
		return errorsmod.Wrapf(ErrInvalidMarketId, "market id cannot be zero")
	}

	// Validate price
	if msg.Price.IsNil() || !msg.Price.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidPrice, "price must be positive")
	}

	// Validate base quantity
	if msg.BaseQuantity.IsNil() || !msg.BaseQuantity.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidQuantity, "base quantity must be positive")
	}

	// Validate order type (must be limit order type)
	if msg.OrderType != OrderType_ORDER_TYPE_LIMIT_BUY && msg.OrderType != OrderType_ORDER_TYPE_LIMIT_SELL {
		return errorsmod.Wrapf(ErrInvalidOrderType, "invalid order type for limit order: %s", msg.OrderType.String())
	}

	return nil
}

func (msg MsgPlaceMarketOrder) ValidateBasic() error {
	// Validate creator address
	if _, err := sdk.AccAddressFromBech32(msg.Creator); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid creator address: %s", err)
	}

	// Validate market ID (must be non-zero)
	if msg.MarketId == 0 {
		return errorsmod.Wrapf(ErrInvalidMarketId, "market id cannot be zero")
	}

	// Validate base quantity
	if msg.BaseQuantity.IsNil() || !msg.BaseQuantity.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidQuantity, "base quantity must be positive")
	}

	// Validate order type (must be market order type)
	if msg.OrderType != OrderType_ORDER_TYPE_MARKET_BUY && msg.OrderType != OrderType_ORDER_TYPE_MARKET_SELL {
		return errorsmod.Wrapf(ErrInvalidOrderType, "invalid order type for market order: %s", msg.OrderType.String())
	}

	return nil
}

func (msg MsgUpdateParams) ValidateBasic() error {
	// Validate authority address
	if _, err := sdk.AccAddressFromBech32(msg.Authority); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid authority address: %s", err)
	}

	// Validate params - delegate to params validation
	if err := msg.Params.Validate(); err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "invalid params: %s", err)
	}

	return nil
}

func (msg MsgLiquidatePositions) ValidateBasic() error {
	_, err := sdk.AccAddressFromBech32(msg.Liquidator)
	if err != nil {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidAddress, "invalid liquidator address (%s)", err)
	}

	// Validate positions array
	if len(msg.Positions) == 0 {
		return errorsmod.Wrap(sdkerrors.ErrInvalidRequest, "positions array cannot be empty")
	}

	// Limit positions to prevent DoS
	const maxPositions = 100
	if len(msg.Positions) > maxPositions {
		return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest, "too many positions: %d (max: %d)",
			len(msg.Positions), maxPositions)
	}

	// Validate each position
	for i, position := range msg.Positions {
		if position.MarketId == 0 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
				"invalid market ID at position %d", i)
		}
		if position.PerpetualId == 0 {
			return errorsmod.Wrapf(sdkerrors.ErrInvalidRequest,
				"invalid perpetual ID at position %d", i)
		}
	}

	return nil
}

func (msg MsgWithdraw) ValidateBasic() error {
	// Validate sender address
	if _, err := sdk.AccAddressFromBech32(msg.Sender); err != nil {
		return errorsmod.Wrapf(ErrInvalidAddress, "invalid sender address: %s", err)
	}

	// Validate coin
	if err := msg.Coin.Validate(); err != nil {
		return errorsmod.Wrapf(ErrInvalidCoin, "invalid coin: %s", err)
	}

	// Check for positive amount
	if msg.Coin.Amount.IsNil() || !msg.Coin.Amount.IsPositive() {
		return errorsmod.Wrapf(ErrInvalidCoin, "coin amount must be positive")
	}

	return nil
}
