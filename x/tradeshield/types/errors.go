package types

// DONTCOVER

import "cosmossdk.io/errors"

// x/tradeshield module sentinel errors
var (
	ErrSample                 = errors.Register(ModuleName, 1100, "sample error")
	ErrSpotOrderNotFound      = errors.Register(ModuleName, 1101, "spot order not found")
	ErrPerpetualOrderNotFound = errors.Register(ModuleName, 1102, "perpetual order not found")
	ErrPriceNotFound          = errors.Register(ModuleName, 1103, "price not found")
	ErrSizeZero               = errors.Register(ModuleName, 1104, "zero order ids ")
	ErrInvalidStatus          = errors.Register(ModuleName, 1105, "invalid status")
	ErrZeroMarketPrice        = errors.Register(ModuleName, 1106, "market price is zero")
	ErrHighTolerance          = errors.Register(ModuleName, 1107, "high tolerance")
)
