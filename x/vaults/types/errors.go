package types

// DONTCOVER

import (
	sdkerrors "cosmossdk.io/errors"
)

// x/vaults module sentinel errors
var (
	ErrInvalidSigner       = sdkerrors.Register(ModuleName, 1100, "expected gov account as only signer for proposal message")
	ErrSample              = sdkerrors.Register(ModuleName, 1101, "sample error")
	ErrVaultNotFound       = sdkerrors.Register(ModuleName, 1102, "vault not found")
	ErrInvalidAction       = sdkerrors.Register(ModuleName, 1103, "invalid action")
	ErrInsufficientShares  = sdkerrors.Register(ModuleName, 1104, "insufficient shares")
	ErrNoShares            = sdkerrors.Register(ModuleName, 1105, "no shares")
	ErrInvalidDepositDenom = sdkerrors.Register(ModuleName, 1106, "invalid deposit denom")
)
