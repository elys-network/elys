package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/estaking module sentinel errors
var (
	ErrNoInflationaryParams = errorsmod.Register(ModuleName, 14, "no inflationary rewards params")
)
