package types

import (
	fmt "fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// TotalCommitmentInfo
// Stores the
type TotalCommitmentInfo struct {
	// Revenue tracking per pool, key => (poolId)
	PoolRevenueTrack map[string]sdk.Dec
}

// Returns the pool revenue tracking key.
// Unique per pool per epoch, clean once complete the calculation.
func GetPoolRevenueTrackKey(poolId uint64) string {
	return fmt.Sprintf("pool_revenue_%d", poolId)
}
