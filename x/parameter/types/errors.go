package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/parameter module sentinel errors
var (
	ErrInvalidMinCommissionRate = sdkerrors.Register(ModuleName, 1101, "invalid min commission rate")
	ErrInvalidMaxVotingPower    = sdkerrors.Register(ModuleName, 1102, "invalid max voting power")
	ErrInvalidMinSelfDelegation = sdkerrors.Register(ModuleName, 1103, "invalid min self delegation")
	ErrInvalidBrokerAddress     = sdkerrors.Register(ModuleName, 1104, "invalid broker address")
)
