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
	ErrInsufficientWithdrawableTokens = errorsmod.Register(ModuleName, 1011, "Once you deposit into the liquidity pool, your funds will be locked for 1 hour. During this time the withdrawal for this specific deposit won't be available")
	ErrVestNowIsNotEnabled            = errorsmod.Register(ModuleName, 1012, "instant vesting is not enabled at this stage")
	ErrUnsupportedWithdrawMode        = errorsmod.Register(ModuleName, 1013, "unsupported withdraw mode")
	ErrUnsupportedUncommitToken       = errorsmod.Register(ModuleName, 1014, "unsupported uncommit token")
	ErrAirdropAlreadyClaimed          = errorsmod.Register(ModuleName, 1015, "airdrop already claimed")
	ErrAirdropNotStarted              = errorsmod.Register(ModuleName, 1016, "airdrop not started")
	ErrAirdropEnded                   = errorsmod.Register(ModuleName, 1017, "airdrop ended")
	ErrMaxElysAmountReached           = errorsmod.Register(ModuleName, 1018, "maximum elys amount reached")
	ErrMaxEdenAmountReached           = errorsmod.Register(ModuleName, 1019, "maximum eden amount reached")
	ErrKolNotFound                    = errorsmod.Register(ModuleName, 1020, "kol not found")
	ErrKolAlreadyClaimed              = errorsmod.Register(ModuleName, 1021, "elys already claimed")
	ErrKolRefunded                    = errorsmod.Register(ModuleName, 1022, "cannot claim elys, refund has been processed")
	ErrClaimNotEnabled                = errorsmod.Register(ModuleName, 1023, "claim not enabled")
	ErrRewardProgramNotFound          = errorsmod.Register(ModuleName, 1024, "reward program not found")
	ErrRewardProgramAlreadyClaimed    = errorsmod.Register(ModuleName, 1025, "reward program already claimed")
	ErrRewardProgramNotStarted        = errorsmod.Register(ModuleName, 1026, "reward program not started")
	ErrRewardProgramEnded             = errorsmod.Register(ModuleName, 1027, "reward program ended")
)
