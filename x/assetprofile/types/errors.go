package types

// DONTCOVER

import "cosmossdk.io/errors"

// x/assetprofile module sentinel errors
var (
	ErrAssetProfileNotFound          = errors.Register(ModuleName, 1, "asset profile not found for denom")
	ErrChannelIdAndDenomHashMismatch = errors.Register(ModuleName, 2, "channel id and denom hash mismatch")
	ErrNotValidIbcDenom              = errors.Register(ModuleName, 3, "not valid ibc denom")
)
