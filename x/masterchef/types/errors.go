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
)
