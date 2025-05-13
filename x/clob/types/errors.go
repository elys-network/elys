package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

var (
	ErrSubAccountNotFound        = errorsmod.Register(ModuleName, 1, "subAccount not found")
	ErrPerpetualOwnerNotFound    = errorsmod.Register(ModuleName, 2, "perpetual owner not found")
	ErrPerpetualNotFound         = errorsmod.Register(ModuleName, 3, "perpetual not found")
	ErrPerpetualMarketNotFound   = errorsmod.Register(ModuleName, 4, "perpetual market not found")
	ErrPerpetualOrderNotFound    = errorsmod.Register(ModuleName, 5, "perpetual order not found")
	ErrFundingRateNotFound       = errorsmod.Register(ModuleName, 6, "funding rate not found")
	ErrInsufficientInsuranceFund = errorsmod.Register(ModuleName, 7, "insufficient insurance fund")
)
