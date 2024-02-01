package types

// DONTCOVER

import (
	"cosmossdk.io/errors"
)

// x/amm module sentinel errors
var (
	ErrLimitMaxAmount      = errors.Register(ModuleName, 6, "calculated amount is larger than max amount")
	ErrLimitMinAmount      = errors.Register(ModuleName, 7, "calculated amount is less than min amount")
	ErrInvalidMathApprox   = errors.Register(ModuleName, 8, "invalid calculated result")
	ErrInvalidPool         = errors.Register(ModuleName, 10, "attempting to create an invalid pool")
	ErrDenomNotFoundInPool = errors.Register(ModuleName, 11, "denom does not exist in pool")

	ErrEmptyRoutes     = errors.Register(ModuleName, 21, "routes not defined")
	ErrNegativeSwapFee = errors.Register(ModuleName, 23, "swap fee is negative")
	ErrNegativeExitFee = errors.Register(ModuleName, 24, "exit fee is negative")
	ErrTooMuchSwapFee  = errors.Register(ModuleName, 25, "swap fee should be less than 1 (100%)")
	ErrTooMuchExitFee  = errors.Register(ModuleName, 26, "exit fee should be less than 1 (100%)")

	ErrTooManyTokensOut = errors.Register(ModuleName, 31, "tx is trying to get more tokens out of the pool than exist")

	ErrInvalidPoolId      = errors.Register(ModuleName, 91, "invalid pool id")
	ErrInvalidSwapMsgType = errors.Register(ModuleName, 92, "unexpected swap message type")

	ErrSameDenom              = errors.Register(ModuleName, 101, "denom in and denom out are the same")
	ErrPoolNotFound           = errors.Register(ModuleName, 102, "pool not found")
	ErrAmountTooLow           = errors.Register(ModuleName, 103, "amount too low")
	ErrInvalidDenom           = errors.Register(ModuleName, 104, "invalid denom")
	ErrInvalidDiscount        = errors.Register(ModuleName, 105, "invalid discount")
	ErrInitialSpotPriceIsZero = errors.Register(ModuleName, 106, "initial spot price is zero")
	ErrSpotPriceIsZero        = errors.Register(ModuleName, 107, "spot price is zero")

	ErrPoolParamsShouldNotBeNil         = errors.Register(ModuleName, 108, "pool params should not be nil")
	ErrSwapFeeShouldNotExceedTwoPercent = errors.Register(ModuleName, 109, "swap fee should not exceed 2%")
	ErrExitFeeShouldNotExceedTwoPercent = errors.Register(ModuleName, 110, "exit fee should not exceed 2%")
	ErrFeeShouldNotBeNegative           = errors.Register(ModuleName, 111, "fee should not be negative")
)

const (
	InvalidInputDenomsErrFormat                  = "input denoms must already exist in the pool (%s)"
	FormatRepeatingPoolAssetsNotAllowedErrFormat = "repeating pool assets not allowed, found %s"
	FormatNoPoolAssetFoundErrFormat              = "can't find the PoolAsset (%s)"
	ErrMsgFormatSharesLargerThanMax              = "%s resulted shares is larger than the max amount of %s"
)
