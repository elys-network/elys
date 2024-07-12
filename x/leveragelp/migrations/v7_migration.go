package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	// Traverse positions and update lp amount and health
	// Liquidate <1.1 positions
	// Update data structure
	positions := m.keeper.GetAllPositions(ctx)
	pools := m.keeper.GetAllPools(ctx)
	for _, pool := range pools {
		m.keeper.DeletePoolPosIdsLiquidationSorted(ctx, pool.AmmPoolId)
		m.keeper.DeletePoolPosIdsStopLossSorted(ctx, pool.AmmPoolId)
	}
	for _, position := range positions {
		m.keeper.SetSortedLiquidationAndStopLoss(ctx, position)
	}

	return nil
}
