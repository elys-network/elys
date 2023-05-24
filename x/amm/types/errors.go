package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/amm module sentinel errors
var (
	ErrInvalidPool = sdkerrors.Register(ModuleName, 10, "attempting to create an invalid pool")

	ErrNegativeSwapFee = sdkerrors.Register(ModuleName, 23, "swap fee is negative")
	ErrNegativeExitFee = sdkerrors.Register(ModuleName, 24, "exit fee is negative")
	ErrTooMuchSwapFee  = sdkerrors.Register(ModuleName, 25, "swap fee should be lesser than 1 (100%)")
	ErrTooMuchExitFee  = sdkerrors.Register(ModuleName, 26, "exit fee should be lesser than 1 (100%)")
)
