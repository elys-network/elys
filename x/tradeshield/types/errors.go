package types

// DONTCOVER

import "cosmossdk.io/errors"

// x/tradeshield module sentinel errors
var (
	ErrSample                 = errors.Register(ModuleName, 1100, "sample error")
	ErrSpotOrderNotFound      = errors.Register(ModuleName, 1101, "spot order not found")
	ErrPerpetualOrderNotFound = errors.Register(ModuleName, 1102, "perpetual order not found")
	ErrPriceNotFound          = errors.Register(ModuleName, 1103, "price not found")
)
