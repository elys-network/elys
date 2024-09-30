package types

// DONTCOVER

import errorsmod "cosmossdk.io/errors"

// x/tradeshield module sentinel errors
var (
	ErrSample        = errorsmod.Register(ModuleName, 1100, "sample error")
	ErrOrderNotFound = errorsmod.Register(ModuleName, 1101, "order not found")
)
