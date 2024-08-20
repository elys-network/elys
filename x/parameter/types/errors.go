package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/parameter module sentinel errors
var (
	ErrInvalidMinCommissionRate    = errorsmod.Register(ModuleName, 1101, "invalid min commission rate")
	ErrInvalidMaxVotingPower       = errorsmod.Register(ModuleName, 1102, "invalid max voting power")
	ErrInvalidMinSelfDelegation    = errorsmod.Register(ModuleName, 1103, "invalid min self delegation")
	ErrInvalidBrokerAddress        = errorsmod.Register(ModuleName, 1104, "invalid broker address")
	ErrInvalidRewardsDataLifecycle = errorsmod.Register(ModuleName, 1105, "invalid rewards data lifecycle")
)
