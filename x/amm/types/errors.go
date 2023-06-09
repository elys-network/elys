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

	ErrInvalidPoolId = sdkerrors.Register(ModuleName, 91, "invalid pool id")
)

const (
	invalidInputDenomsErrFormat = "input denoms must already exist in the pool (%s)"

	formatRepeatingPoolAssetsNotAllowedErrFormat = "repeating pool assets not allowed, found %s"
)
