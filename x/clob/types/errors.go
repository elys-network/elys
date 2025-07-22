package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

var (
	// Object not found errors
	ErrSubAccountNotFound      = errorsmod.Register(ModuleName, 1, "subAccount not found")
	ErrPerpetualOwnerNotFound  = errorsmod.Register(ModuleName, 2, "perpetual owner not found")
	ErrPerpetualNotFound       = errorsmod.Register(ModuleName, 3, "perpetual not found")
	ErrPerpetualMarketNotFound = errorsmod.Register(ModuleName, 4, "perpetual market not found")
	ErrPerpetualOrderNotFound  = errorsmod.Register(ModuleName, 5, "perpetual order not found")
	ErrFundingRateNotFound     = errorsmod.Register(ModuleName, 6, "funding rate not found")
	ErrPerpetualADLNotFound    = errorsmod.Register(ModuleName, 8, "perpetual ADL not found")

	// Validation errors
	ErrInvalidPrice     = errorsmod.Register(ModuleName, 9, "invalid price")
	ErrInvalidQuantity  = errorsmod.Register(ModuleName, 10, "invalid quantity")
	ErrInvalidOrderType = errorsmod.Register(ModuleName, 11, "invalid order type")
	ErrInvalidAddress   = errorsmod.Register(ModuleName, 12, "invalid address")
	ErrInvalidCoin      = errorsmod.Register(ModuleName, 13, "invalid coin")
	ErrInvalidMarketId  = errorsmod.Register(ModuleName, 14, "invalid market id")
	ErrInvalidTradeSize = errorsmod.Register(ModuleName, 15, "invalid trade size")

	// Business logic errors
	ErrInsufficientInsuranceFund = errorsmod.Register(ModuleName, 7, "insufficient insurance fund")
	ErrInsufficientBalance       = errorsmod.Register(ModuleName, 17, "insufficient balance")
	ErrInsufficientMargin        = errorsmod.Register(ModuleName, 18, "insufficient margin")
	ErrOrderNotFilled            = errorsmod.Register(ModuleName, 19, "order cannot be fully filled")
	ErrMarketNotActive           = errorsmod.Register(ModuleName, 20, "market is not active")
	ErrPositionNotHealthy        = errorsmod.Register(ModuleName, 21, "position is not healthy")

	// Order-specific errors
	ErrNotBuyOrder        = errorsmod.Register(ModuleName, 22, "order is not a buy order")
	ErrNotSellOrder       = errorsmod.Register(ModuleName, 23, "order is not a sell order")
	ErrOrderAlreadyExists = errorsmod.Register(ModuleName, 24, "order already exists")
	ErrMaxOrdersExceeded  = errorsmod.Register(ModuleName, 25, "maximum orders exceeded")

	// Calculation errors
	ErrDivisionByZero      = errorsmod.Register(ModuleName, 26, "division by zero")
	ErrCalculationOverflow = errorsmod.Register(ModuleName, 27, "calculation overflow")
	ErrNegativeResult      = errorsmod.Register(ModuleName, 28, "calculation resulted in negative value")

	// System errors
	ErrImplementationIncomplete = errorsmod.Register(ModuleName, 16, "implementation incomplete")
	ErrUnmarshalFailed          = errorsmod.Register(ModuleName, 29, "failed to unmarshal data")
	ErrMarshalFailed            = errorsmod.Register(ModuleName, 30, "failed to marshal data")
	ErrInvalidState             = errorsmod.Register(ModuleName, 31, "invalid state")

	// Rate limiting errors
	ErrRateLimitExceeded       = errorsmod.Register(ModuleName, 32, "rate limit exceeded")
	ErrCircuitBreakerTriggered = errorsmod.Register(ModuleName, 33, "circuit breaker triggered")
)
