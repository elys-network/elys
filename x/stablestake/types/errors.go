package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/stablestake module sentinel errors
var (
	ErrInvalidDepositDenom  = errorsmod.Register(ModuleName, 1, "invalid deposit denom")
	ErrInvalidBorrowDenom   = errorsmod.Register(ModuleName, 2, "invalid borrow denom")
	ErrRedemptionRateIsZero = errorsmod.Register(ModuleName, 3, "redemption rate is zero")
	ErrNegativeBorrowed     = errorsmod.Register(ModuleName, 4, "negative borrowed amount")
	ErrMaxBorrowAmount      = errorsmod.Register(ModuleName, 5, "cannot borrow more than 90% of total stablestake pool.")
	ErrInvalidParams        = errorsmod.Register(ModuleName, 6, "invalid params")
	ErrInvalidWithdraw      = errorsmod.Register(ModuleName, 7, "cannot withdraw, max borrow ratio limit reached")
	ErrPoolNotFound         = errorsmod.Register(ModuleName, 8, "pool not found")
	ErrPoolAlreadyExists    = errorsmod.Register(ModuleName, 9, "pool already exists")
)
