package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/stablestake module sentinel errors
var (
	ErrInvalidDepositDenom  = sdkerrors.Register(ModuleName, 1, "invalid deposit denom")
	ErrInvalidBorrowDenom   = sdkerrors.Register(ModuleName, 2, "invalid borrow denom")
	ErrRedemptionRateIsZero = sdkerrors.Register(ModuleName, 3, "redemption rate is zero")
)
