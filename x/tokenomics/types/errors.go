package types

// DONTCOVER

import "cosmossdk.io/errors"

// x/tokenomics module sentinel errors
var (
	ErrAirdropExpired = errors.Register(ModuleName, 1, "Expired airdrop")
)
