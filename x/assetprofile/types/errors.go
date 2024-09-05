package types

// DONTCOVER

import "cosmossdk.io/errors"

// x/assetprofile module sentinel errors
var (
	ErrAssetProfileNotFound          = errors.Register(ModuleName, 1, "asset profile not found for denom")
	ErrChannelIdAndDenomHashMismatch = errors.Register(ModuleName, 2, "channel id and denom hash mismatch")
	ErrNotValidIbcDenom              = errors.Register(ModuleName, 3, "not valid ibc denom")
	ErrDecimalsInvalid               = errors.Register(ModuleName, 4, "error decimals have to be a value between 6 and 18")
	ErrInvalidBaseDenom              = errors.Register(ModuleName, 5, "error invalid base denom")
)
