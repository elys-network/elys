package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TotalCommitmentInfo
// Stores the
type TotalCommitmentInfo struct {
	// Total Elys staked
	TotalElysBonded sdk.Int
	// Total LP committed
	TotalLPCommitted sdk.Int
	// Total Eden committed
	TotalEdenCommitted sdk.Int
	// Total Eden boost committed
	TotalEdenBoostCommitted sdk.Int
}
