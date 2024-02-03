package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/burner module sentinel errors
var (
	ErrSample = errorsmod.Register(ModuleName, 1100, "sample error")
)
