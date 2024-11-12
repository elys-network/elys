package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/tier module sentinel errors
var (
	ErrSample   = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrNotFound = errorsmod.Register(ModuleName, 1101, "asset not found in asset profile")
)
