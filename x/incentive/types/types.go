package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
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

// Returns the rewards wallet per pool
func GetLPRewardsPoolAddress(poolId uint64) sdk.AccAddress {
	return authtypes.NewModuleAddress(fmt.Sprintf("lp_rewards_pool_%d", poolId))
}
