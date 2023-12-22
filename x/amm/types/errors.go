package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/amm module sentinel errors
var (
	ErrLimitMaxAmount      = sdkerrors.Register(ModuleName, 6, "calculated amount is larger than max amount")
	ErrLimitMinAmount      = sdkerrors.Register(ModuleName, 7, "calculated amount is lesser than min amount")
	ErrInvalidMathApprox   = sdkerrors.Register(ModuleName, 8, "invalid calculated result")
	ErrInvalidPool         = sdkerrors.Register(ModuleName, 10, "attempting to create an invalid pool")
	ErrDenomNotFoundInPool = sdkerrors.Register(ModuleName, 11, "denom does not exist in pool")

	ErrEmptyRoutes     = sdkerrors.Register(ModuleName, 21, "routes not defined")
	ErrNegativeSwapFee = sdkerrors.Register(ModuleName, 23, "swap fee is negative")
	ErrNegativeExitFee = sdkerrors.Register(ModuleName, 24, "exit fee is negative")
	ErrTooMuchSwapFee  = sdkerrors.Register(ModuleName, 25, "swap fee should be lesser than 1 (100%)")
	ErrTooMuchExitFee  = sdkerrors.Register(ModuleName, 26, "exit fee should be lesser than 1 (100%)")

	ErrTooManyTokensOut = sdkerrors.Register(ModuleName, 31, "tx is trying to get more tokens out of the pool than exist")

	ErrInvalidPoolId      = sdkerrors.Register(ModuleName, 91, "invalid pool id")
	ErrInvalidSwapMsgType = sdkerrors.Register(ModuleName, 92, "unexpected swap message type")

	ErrSameDenom              = sdkerrors.Register(ModuleName, 101, "denom in and denom out are the same")
	ErrPoolNotFound           = sdkerrors.Register(ModuleName, 102, "pool not found")
	ErrAmountTooLow           = sdkerrors.Register(ModuleName, 103, "amount too low")
	ErrInvalidDenom           = sdkerrors.Register(ModuleName, 104, "invalid denom")
	ErrInvalidDiscount        = sdkerrors.Register(ModuleName, 105, "invalid discount")
	ErrInitialSpotPriceIsZero = sdkerrors.Register(ModuleName, 106, "initial spot price is zero")
	ErrSpotPriceIsZero        = sdkerrors.Register(ModuleName, 107, "spot price is zero")
)

const (
	InvalidInputDenomsErrFormat                  = "input denoms must already exist in the pool (%s)"
	FormatRepeatingPoolAssetsNotAllowedErrFormat = "repeating pool assets not allowed, found %s"
	FormatNoPoolAssetFoundErrFormat              = "can't find the PoolAsset (%s)"
	ErrMsgFormatSharesLargerThanMax              = "%s resulted shares is larger than the max amount of %s"
)
