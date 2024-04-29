package types

import (
	errorsmod "cosmossdk.io/errors"
)

// DONTCOVER

// x/commitment module sentinel errors
var (
	ErrCommitmentsNotFound            = errorsmod.Register(ModuleName, 1001, "commitments not found for creator")
	ErrInsufficientRewardsUnclaimed   = errorsmod.Register(ModuleName, 1002, "insufficient unclaimed tokens for creator and denom")
	ErrInsufficientCommittedTokens    = errorsmod.Register(ModuleName, 1003, "insufficient committed tokens for creator and denom")
	ErrInsufficientClaimed            = errorsmod.Register(ModuleName, 1004, "insufficient claimed tokens for creator and denom")
	ErrInvalidAmount                  = errorsmod.Register(ModuleName, 1005, "invalid amount")
	ErrInvalidDenom                   = errorsmod.Register(ModuleName, 1006, "invalid denom")
	ErrInsufficientVestingTokens      = errorsmod.Register(ModuleName, 1007, "insufficient vesting tokens for creator and denom")
	ErrCommitDisabled                 = errorsmod.Register(ModuleName, 1008, "commitment disabled for denom")
	ErrWithdrawDisabled               = errorsmod.Register(ModuleName, 1009, "withdraw disabled for denom")
	ErrExceedMaxVestings              = errorsmod.Register(ModuleName, 1010, "exceed maximum allowed vestings")
	ErrInsufficientWithdrawableTokens = errorsmod.Register(ModuleName, 1011, "insufficient withdrawable tokens")
	ErrVestNowIsNotEnabled            = errorsmod.Register(ModuleName, 1012, "instant vesting is not enabled at this stage")
	ErrUnsupportedWithdrawMode        = errorsmod.Register(ModuleName, 1013, "unsupported withdraw mode")
	ErrUnsupportedUncommitToken       = errorsmod.Register(ModuleName, 1014, "unsupported uncommit token")
)
