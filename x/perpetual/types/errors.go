package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/perpetual module sentinel errors
var (
	ErrMTPDoesNotExist        = errorsmod.Register(ModuleName, 1, "mtp not found")
	ErrMTPInvalid             = errorsmod.Register(ModuleName, 2, "mtp invalid")
	ErrMTPDisabled            = errorsmod.Register(ModuleName, 3, "perpetual not enabled for pool")
	ErrUnknownRequest         = errorsmod.Register(ModuleName, 4, "unknown request")
	ErrMTPHealthy             = errorsmod.Register(ModuleName, 5, "mtp health above force close threshold")
	ErrInvalidPosition        = errorsmod.Register(ModuleName, 6, "mtp position invalid")
	ErrMaxOpenPositions       = errorsmod.Register(ModuleName, 7, "max open positions reached")
	ErrUnauthorised           = errorsmod.Register(ModuleName, 8, "address not on whitelist")
	ErrBorrowTooLow           = errorsmod.Register(ModuleName, 9, "borrowed amount is too low")
	ErrBorrowTooHigh          = errorsmod.Register(ModuleName, 10, "borrowed amount is higher than pool depth")
	ErrCustodyTooHigh         = errorsmod.Register(ModuleName, 11, "custody amount is higher than pool depth")
	ErrMTPUnhealthy           = errorsmod.Register(ModuleName, 12, "mtp health would be too low for safety factor")
	ErrInvalidCollateralAsset = errorsmod.Register(ModuleName, 13, "invalid collateral asset")
	ErrInvalidBorrowingAsset  = errorsmod.Register(ModuleName, 14, "invalid borrowing asset")
	ErrPoolDoesNotExist       = errorsmod.Register(ModuleName, 15, "perpetual pool does not exist")
	ErrBalanceNotAvailable    = errorsmod.Register(ModuleName, 18, "user does not have enough balance of the required coin")
	ErrAmountTooLow           = errorsmod.Register(ModuleName, 32, "Tx amount is too low")
	ErrPerpetualDisabled      = errorsmod.Register(ModuleName, 33, "perpetual disabled pool")
	ErrDenomNotFound          = errorsmod.Register(ModuleName, 34, "denom not found")
	ErrInvalidLeverage        = errorsmod.Register(ModuleName, 35, "invalid leverage")
	ErrInvalidCloseSize       = errorsmod.Register(ModuleName, 36, "invalid close size")
	ErrCalcMinCollateral      = errorsmod.Register(ModuleName, 37, "error calculating min collateral")
	ErrInvalidTakeProfitPrice = errorsmod.Register(ModuleName, 38, "error invalid profit price ")
	ErrInvalidTradingAsset    = errorsmod.Register(ModuleName, 39, "invalid trading asset")
	ErrInvalidAmount          = errorsmod.Register(ModuleName, 40, "invalid amount")
	ErrInvalidPrice           = errorsmod.Register(ModuleName, 41, "error invalid price ")
	ErrPoolHasToBeOracle      = errorsmod.Register(ModuleName, 42, "pool has to be oracle enabled")
	ErrZeroCustodyAmount      = errorsmod.Register(ModuleName, 43, "Custody amount is zero")
	ErrPoolNotEnabled         = errorsmod.Register(ModuleName, 44, "pool is not enabled")
)
