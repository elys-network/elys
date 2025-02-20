package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	// reset params
	// legacy := m.keeper.GetLegacyParams(ctx)
	// params := types.NewParams(
	// 	legacy.MinCommissionRate,           // min commission 0.05
	// 	legacy.MaxVotingPower,              // max voting power 0.66
	// 	legacy.MinSelfDelegation,           // min self delegation
	// 	uint64(legacy.TotalBlocksPerYear),  // total blocks per year
	// 	uint64(legacy.RewardsDataLifetime), // 24 hrs
	// )
	// m.keeper.SetParams(ctx, params)

	return nil
}
