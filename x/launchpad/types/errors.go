package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/launchpad module sentinel errors
var (
	ErrNotEnabledSpendingToken  = errorsmod.Register(ModuleName, 1, "not enabled spending token")
	ErrLaunchpadNotStarted      = errorsmod.Register(ModuleName, 2, "launchpad not started")
	ErrLaunchpadAlreadyFinished = errorsmod.Register(ModuleName, 3, "launchpad already finished")
	ErrOverflowTotalReserve     = errorsmod.Register(ModuleName, 4, "buy amount overflow total reserve")
)
