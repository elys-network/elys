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
	ErrPerpetualADLNotFound      = errorsmod.Register(ModuleName, 8, "perpetual ADL not found")
	ErrInvalidPrice              = errorsmod.Register(ModuleName, 9, "invalid price")
	ErrInvalidQuantity           = errorsmod.Register(ModuleName, 10, "invalid quantity")
	ErrInvalidOrderType          = errorsmod.Register(ModuleName, 11, "invalid order type")
	ErrInvalidAddress            = errorsmod.Register(ModuleName, 12, "invalid address")
	ErrInvalidCoin               = errorsmod.Register(ModuleName, 13, "invalid coin")
	ErrInvalidMarketId           = errorsmod.Register(ModuleName, 14, "invalid market id")
	ErrInvalidTradeSize          = errorsmod.Register(ModuleName, 15, "invalid trade size")
	ErrImplementationIncomplete  = errorsmod.Register(ModuleName, 16, "implementation incomplete")
)
