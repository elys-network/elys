package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/leveragelp module sentinel errors
var (
	ErrMTPDoesNotExist        = sdkerrors.Register(ModuleName, 1, "mtp not found")
	ErrMTPInvalid             = sdkerrors.Register(ModuleName, 2, "mtp invalid")
	ErrMTPDisabled            = sdkerrors.Register(ModuleName, 3, "leveragelp not enabled for pool")
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
	ErrLeveragelpDisabled     = sdkerrors.Register(ModuleName, 33, "leveragelp disabled pool")
	ErrAmmPoolNotFound        = sdkerrors.Register(ModuleName, 34, "amm pool not found")
)
