package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/tier module sentinel errors
var (
	ErrSample   = sdkerrors.Register(ModuleName, 1100, "sample error")
	ErrNotFound = sdkerrors.Register(ModuleName, 1101, "asset not found in asset profiler")
)
