package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	// Traverse positions and update lp amount and health
	// Update data structure
	pools := m.keeper.GetAllPools(ctx)
	for _, pool := range pools {
		m.keeper.DeletePoolPosIdsLiquidationSorted(ctx, pool.AmmPoolId)
		m.keeper.DeletePoolPosIdsStopLossSorted(ctx, pool.AmmPoolId)
	}
	openCount := uint64(0)

	m.keeper.SetOpenPositionCount(ctx, openCount)
	return nil
}
