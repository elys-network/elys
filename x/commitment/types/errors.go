package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/commitment module sentinel errors
var (
	ErrCommitmentsNotFound            = sdkerrors.Register(ModuleName, 1001, "commitments not found for creator")
	ErrInsufficientRewardsUnclaimed   = sdkerrors.Register(ModuleName, 1002, "insufficient unclaimed tokens for creator and denom")
	ErrInsufficientCommittedTokens    = sdkerrors.Register(ModuleName, 1003, "insufficient committed tokens for creator and denom")
	ErrInvalidAmount                  = sdkerrors.Register(ModuleName, 1004, "invalid amount")
	ErrInvalidDenom                   = sdkerrors.Register(ModuleName, 1005, "invalid denom")
	ErrInsufficientVestingTokens      = sdkerrors.Register(ModuleName, 1006, "insufficient vesting tokens for creator and denom")
	ErrCommitDisabled                 = sdkerrors.Register(ModuleName, 1007, "commitment disabled for denom")
	ErrWithdrawDisabled               = sdkerrors.Register(ModuleName, 1008, "withdraw disabled for denom")
	ErrExceedMaxVestings              = sdkerrors.Register(ModuleName, 1009, "exceed maximum allowed vestings")
	ErrInsufficientWithdrawableTokens = sdkerrors.Register(ModuleName, 1010, "insufficient withdrawable tokens")
	ErrVestNowIsNotEnabled            = sdkerrors.Register(ModuleName, 1011, "instant vesting is not enabled at this stage")
)
