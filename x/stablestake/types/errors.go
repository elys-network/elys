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
)
