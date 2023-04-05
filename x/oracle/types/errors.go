package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/oracle module sentinel errors
var (
	ErrNotAvailable         = sdkerrors.Register(ModuleName, 1500, "sample error")
	ErrInvalidPacketTimeout = sdkerrors.Register(ModuleName, 1501, "invalid packet timeout")
	ErrInvalidVersion       = sdkerrors.Register(ModuleName, 1502, "invalid version")
	ErrNotAPriceFeeder      = sdkerrors.Register(ModuleName, 1503, "not a price feeder")
	ErrPriceFeederNotActive = sdkerrors.Register(ModuleName, 1504, "price feeder is not active")
	ErrNotModuleAdmin       = sdkerrors.Register(ModuleName, 1505, "not a module admin")
)
