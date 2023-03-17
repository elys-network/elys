package types

// DONTCOVER

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// x/commitment module sentinel errors
var (
	ErrCommitmentsNotFound           = sdkerrors.Register(ModuleName, 1, "commitments not found for creator")
	ErrInsufficientUncommittedTokens = sdkerrors.Register(ModuleName, 2, "insufficient uncommitted tokens for creator and denom")
	ErrInsufficientCommittedTokens   = sdkerrors.Register(ModuleName, 3, "insufficient committed tokens for creator and denom")
	ErrInvalidAmount                 = sdkerrors.Register(ModuleName, 4, "invalid amount")
)
