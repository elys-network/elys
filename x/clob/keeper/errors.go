package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/clob/types"
)

// Enhanced error messages with context

// WrapInsufficientBalanceError wraps insufficient balance errors with context
func WrapInsufficientBalanceError(available, required math.Int, denom string, operation string) error {
	return errorsmod.Wrapf(types.ErrInsufficientBalance,
		"insufficient balance for %s: available %s %s, required %s %s",
		operation, available, denom, required, denom)
}

// WrapInsufficientLiquidityError wraps insufficient liquidity errors with context
func WrapInsufficientLiquidityError(available, required math.LegacyDec, orderType string) error {
	return errorsmod.Wrapf(types.ErrInsufficientLiquidity,
		"insufficient liquidity for %s order: available %s, required %s",
		orderType, available, required)
}

// WrapPriceNotFoundError wraps price not found errors with context
func WrapPriceNotFoundError(denom string, source string) error {
	return errorsmod.Wrapf(types.ErrPriceNotFound,
		"%s price not found for %s", source, denom)
}

// WrapInvalidPriceError wraps invalid price errors with context
func WrapInvalidPriceError(price interface{}, reason string) error {
	return errorsmod.Wrapf(types.ErrInvalidPrice,
		"invalid price %v: %s", price, reason)
}

// WrapMarketNotFoundError wraps market not found errors with context
func WrapMarketNotFoundError(marketId uint64) error {
	return errorsmod.Wrapf(types.ErrPerpetualMarketNotFound,
		"perpetual market with id %d not found", marketId)
}

// WrapOrderNotFoundError wraps order not found errors with context
func WrapOrderNotFoundError(orderId uint64, marketId uint64) error {
	return errorsmod.Wrapf(types.ErrOrderNotFound,
		"order %d not found in market %d", orderId, marketId)
}

// WrapPositionNotFoundError wraps position not found errors with context
func WrapPositionNotFoundError(owner string, perpetualId uint64) error {
	return errorsmod.Wrapf(types.ErrPerpetualNotFound,
		"position %d not found for owner %s", perpetualId, owner)
}

// WrapInvalidQuantityError wraps invalid quantity errors with context
func WrapInvalidQuantityError(quantity math.LegacyDec, reason string) error {
	return errorsmod.Wrapf(types.ErrInvalidQuantity,
		"invalid quantity %s: %s", quantity, reason)
}

// WrapLiquidationError wraps liquidation errors with context
func WrapLiquidationError(perpetualId uint64, reason string) error {
	return errorsmod.Wrapf(types.ErrLiquidation,
		"liquidation failed for position %d: %s", perpetualId, reason)
}

// WrapMarginCheckError wraps margin check errors with context
func WrapMarginCheckError(required, available math.Int, operation string) error {
	return errorsmod.Wrapf(types.ErrInsufficientMargin,
		"margin check failed for %s: required %s, available %s",
		operation, required, available)
}

// WrapCalculationError wraps calculation errors with context
func WrapCalculationError(operation string, details string) error {
	return fmt.Errorf("calculation error in %s: %s", operation, details)
}

// WrapValidationError wraps validation errors with context
func WrapValidationError(field string, value interface{}, constraint string) error {
	return fmt.Errorf("validation failed for %s with value %v: %s", field, value, constraint)
}

// WrapSubAccountError wraps sub account errors with context
func WrapSubAccountError(owner string, subAccountId uint64, operation string, err error) error {
	return errorsmod.Wrapf(err,
		"sub account error for owner %s, sub account %d during %s",
		owner, subAccountId, operation)
}

// WrapTransferError wraps transfer errors with context
func WrapTransferError(from, to string, amount sdk.Coins, err error) error {
	return errorsmod.Wrapf(err,
		"transfer failed from %s to %s for amount %s",
		from, to, amount)
}
