package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/assetprofile module sentinel errors
var (
	ErrAssetProfileNotFound = sdkerrors.Register(ModuleName, 1100, "asset profile not found for denom")
)
