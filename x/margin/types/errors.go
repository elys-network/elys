package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/margin module sentinel errors
var (
	ErrMTPDoesNotExist        = sdkerrors.Register(ModuleName, 1, "mtp not found")
	ErrMTPInvalid             = sdkerrors.Register(ModuleName, 2, "mtp invalid")
	ErrMTPDisabled            = sdkerrors.Register(ModuleName, 3, "margin not enabled for pool")
	ErrUnknownRequest         = sdkerrors.Register(ModuleName, 4, "unknown request")
	ErrMTPHealthy             = sdkerrors.Register(ModuleName, 5, "mtp health above force close threshold")
	ErrInvalidPosition        = sdkerrors.Register(ModuleName, 6, "mtp position invalid")
	ErrMaxOpenPositions       = sdkerrors.Register(ModuleName, 7, "max open positions reached")
	ErrUnauthorised           = sdkerrors.Register(ModuleName, 8, "address not on whitelist")
	ErrBorrowTooLow           = sdkerrors.Register(ModuleName, 9, "borrowed amount is too low")
	ErrBorrowTooHigh          = sdkerrors.Register(ModuleName, 10, "borrowed amount is higher than pool depth")
	ErrCustodyTooHigh         = sdkerrors.Register(ModuleName, 11, "custody amount is higher than pool depth")
	ErrMTPUnhealthy           = sdkerrors.Register(ModuleName, 12, "mtp health would be too low for safety factor")
	ErrInvalidCollateralAsset = sdkerrors.Register(ModuleName, 13, "invalid collateral asset")
	ErrInvalidBorrowingAsset  = sdkerrors.Register(ModuleName, 14, "invalid borrowing asset")
	ErrPoolDoesNotExist       = sdkerrors.Register(ModuleName, 15, "pool does not exist")
	ErrBalanceNotAvailable    = sdkerrors.Register(ModuleName, 18, "user does not have enough balance of the required coin")
	ErrAmountTooLow           = sdkerrors.Register(ModuleName, 32, "Tx amount is too low")
	ErrMarginDisabled         = sdkerrors.Register(ModuleName, 33, "margin disabled pool")
	ErrDenomNotFound          = sdkerrors.Register(ModuleName, 34, "denom not found")
	ErrInvalidLeverage        = sdkerrors.Register(ModuleName, 35, "invalid leverage")
	ErrInvalidCloseSize       = sdkerrors.Register(ModuleName, 36, "invalid close size")
	ErrCalcMinCollateral      = sdkerrors.Register(ModuleName, 37, "error calculating min collateral")
)
