package types

// DONTCOVER

import (
	errorsmod "cosmossdk.io/errors"
)

// x/launchpad module sentinel errors
var (
	ErrNotEnabledSpendingToken       = errorsmod.Register(ModuleName, 1, "not enabled spending token")
	ErrLaunchpadNotStarted           = errorsmod.Register(ModuleName, 2, "launchpad not started")
	ErrLaunchpadAlreadyFinished      = errorsmod.Register(ModuleName, 3, "launchpad already finished")
	ErrOverflowTotalReserve          = errorsmod.Register(ModuleName, 4, "buy amount overflow total reserve")
	ErrLaunchpadNotFinished          = errorsmod.Register(ModuleName, 5, "launchpad not finished")
	ErrLaunchpadReturnPeriodFinished = errorsmod.Register(ModuleName, 6, "launchpad return period finished")
	ErrPurchaseOrderNotFound         = errorsmod.Register(ModuleName, 7, "purchase order not found")
	ErrExceedMaxReturnAmount         = errorsmod.Register(ModuleName, 8, "exceed max return amount")
	ErrExceedMaxWithdrawableAmount   = errorsmod.Register(ModuleName, 9, "exceed max withdrawable amount")
	ErrInvalidWithrawAccount         = errorsmod.Register(ModuleName, 10, "invalid withdraw account")
)
