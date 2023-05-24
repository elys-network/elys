package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TotalCommitmentInfo
// Stores the
type TotalCommitmentInfo struct {
	// Total Elys staked
	TotalElysBonded sdk.Int
	// Total Eden + Eden boost committed
	TotalEdenEdenBoostCommitted sdk.Int
	// Gas fees collected and DEX revenus
	TotalFeesCollected sdk.Coins
	// Total Lp Token committed
	TotalLpTokensCommitted map[string]sdk.Int
}
