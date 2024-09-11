package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/masterchef module sentinel errors
var (
	ErrNoInflationaryParams     = errorsmod.Register(ModuleName, 1, "no inflationary rewards params")
	ErrPoolRewardsAccumNotFound = errorsmod.Register(ModuleName, 2, "pool rewards accumulation not found")
	ErrPoolNotFound             = errorsmod.Register(ModuleName, 3, "pool not found")
	ErrInvalidBlockRange        = errorsmod.Register(ModuleName, 4, "invalid block range")
	ErrInvalidAmountPerBlock    = errorsmod.Register(ModuleName, 5, "error invalid amount per block")
	ErrInvalidMinAmount         = errorsmod.Register(ModuleName, 6, "error invalid min amount")
	ErrInvalidPoolMultiplier    = errorsmod.Register(ModuleName, 7, "error invalid pool multiplier")
)
