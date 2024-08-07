package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/oracle module sentinel errors
var (
	ErrNotAvailable         = errorsmod.Register(ModuleName, 1500, "sample error")
	ErrInvalidPacketTimeout = errorsmod.Register(ModuleName, 1501, "invalid packet timeout")
	ErrInvalidVersion       = errorsmod.Register(ModuleName, 1502, "invalid version")
	ErrNotAPriceFeeder      = errorsmod.Register(ModuleName, 1503, "not a price feeder")
	ErrPriceFeederNotActive = errorsmod.Register(ModuleName, 1504, "price feeder is not active")
	ErrNotModuleAdmin       = errorsmod.Register(ModuleName, 1505, "not a module admin")
	ErrAssetWasCreated      = errorsmod.Register(ModuleName, 1506, "asset already exists")
)
