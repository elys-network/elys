package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/accountedpool module sentinel errors
var (
	ErrPoolDoesNotExist          = sdkerrors.Register(ModuleName, 1, "pool does not exist")
	ErrDuplicatedAccountedPoolId = sdkerrors.Register(ModuleName, 2, "duplicated poolId for accountedPool")
)
