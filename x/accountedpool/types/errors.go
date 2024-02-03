package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/accountedpool module sentinel errors
var (
	ErrPoolDoesNotExist          = errorsmod.Register(ModuleName, 1, "pool does not exist")
	ErrDuplicatedAccountedPoolId = errorsmod.Register(ModuleName, 2, "duplicated poolId for accountedPool")
)
