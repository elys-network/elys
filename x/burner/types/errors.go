package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/burner module sentinel errors
var (
	ErrInvalidEpochIdentifier = errorsmod.Register(ModuleName, 1, "invalid epoch identifier")
	ErrInvalidParams          = errorsmod.Register(ModuleName, 2, "invalid param")
)
