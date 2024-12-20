package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/leveragelp module sentinel errors
var (
	ErrPositionDoesNotExist      = errorsmod.Register(ModuleName, 1, "position not found")
	ErrPositionInvalid           = errorsmod.Register(ModuleName, 2, "position invalid")
	ErrPositionDisabled          = errorsmod.Register(ModuleName, 3, "leveragelp not enabled for pool")
	ErrInvalidPosition           = errorsmod.Register(ModuleName, 6, "position position invalid")
	ErrMaxOpenPositions          = errorsmod.Register(ModuleName, 7, "max open positions reached")
	ErrUnauthorised              = errorsmod.Register(ModuleName, 8, "address not on whitelist")
	ErrBorrowTooLow              = errorsmod.Register(ModuleName, 9, "borrowed amount is too low")
	ErrBorrowTooHigh             = errorsmod.Register(ModuleName, 10, "borrowed amount is higher than pool depth")
	ErrCustodyTooHigh            = errorsmod.Register(ModuleName, 11, "custody amount is higher than pool depth")
	ErrPositionUnhealthy         = errorsmod.Register(ModuleName, 12, "position health would be too low for safety factor")
	ErrInvalidCollateralAsset    = errorsmod.Register(ModuleName, 13, "invalid collateral asset")
	ErrInvalidBorrowingAsset     = errorsmod.Register(ModuleName, 14, "invalid borrowing asset")
	ErrPoolDoesNotExist          = errorsmod.Register(ModuleName, 15, "pool does not exist")
	ErrBalanceNotAvailable       = errorsmod.Register(ModuleName, 18, "user does not have enough balance of the required coin")
	ErrAmountTooLow              = errorsmod.Register(ModuleName, 32, "Tx amount is too low")
	ErrLeveragelpDisabled        = errorsmod.Register(ModuleName, 33, "leveragelp disabled pool")
	ErrAmmPoolNotFound           = errorsmod.Register(ModuleName, 34, "amm pool not found")
	ErrOnlyBaseCurrencyAllowed   = errorsmod.Register(ModuleName, 35, "only base currency is allowed for leverage lp")
	ErrInsufficientUsdcAfterOp   = errorsmod.Register(ModuleName, 36, "insufficient amount of usdc after the operation for leveragelp withdrawal")
	ErrInvalidCloseSize          = errorsmod.Register(ModuleName, 37, "invalid close size")
	ErrNegUserAmountAfterRepay   = errorsmod.Register(ModuleName, 38, "negative user amount after repay")
	ErrInvalidLeverage           = errorsmod.Register(ModuleName, 39, "leverage should be same as existing position")
	ErrInvalidCollateral         = errorsmod.Register(ModuleName, 40, "collateral should not be more than total liability")
	ErrPoolLeverageAmountNotZero = errorsmod.Register(ModuleName, 41, "pool leverage amount is greater than zero")
	ErrLeverageTooSmall          = errorsmod.Register(ModuleName, 42, "leverage should be more than or equal to 1")
	ErrMaxLeverageLpExists       = errorsmod.Register(ModuleName, 43, "pool is already leveraged at maximum value")
	ErrUnbondingPoolHealth       = errorsmod.Register(ModuleName, 44, "pool health too low to unbond")
)
