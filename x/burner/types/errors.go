package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/burner module sentinel errors
var (
	ErrInvalidEpochIdentifier = errorsmod.Register(ModuleName, 1, "error invalid epoch identifier")
)
