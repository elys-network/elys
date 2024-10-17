package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/accountedpool module sentinel errors
var (
	ErrPoolDoesNotExist          = errorsmod.Register(ModuleName, 1, "accounted pool does not exist")
	ErrPoolAlreadyExist          = errorsmod.Register(ModuleName, 2, "accounted pool already exist")
	ErrDuplicatedAccountedPoolId = errorsmod.Register(ModuleName, 3, "duplicated poolId for accountedPool")
)
