package migrations

import (
	"fmt"

	"cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	// Traverse positions and update lp amount and health
	// Update data structure
	positions := m.keeper.GetAllPositions(ctx)
	pools := m.keeper.GetAllPools(ctx)
	for _, pool := range pools {
		m.keeper.DeletePoolPosIdsLiquidationSorted(ctx, pool.AmmPoolId)
		m.keeper.DeletePoolPosIdsStopLossSorted(ctx, pool.AmmPoolId)
	}
	openCount := uint64(0)
	for _, position := range positions {
		m.keeper.SetSortedLiquidationAndStopLoss(ctx, position)
		openCount++
	}

	m.keeper.SetOpenPositionCount(ctx, openCount)

	// Liquidate <1.1 positions
	// Q: What will happen if there won't be enough liquidity to return to users(as health for some positions must be below 1) ? Do we need to fill the pool ?
	for _, pool := range pools {
		ammPool, err := m.keeper.GetAmmPool(ctx, pool.AmmPoolId)
		if err != nil {
			ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
			continue
		}
		m.keeper.IteratePoolPosIdsLiquidationSorted(ctx, pool.AmmPoolId, func(posId types.AddressId) bool {
			position, err := m.keeper.GetLegacyPosition(ctx, posId.Address, posId.Id)
			if err != nil {
				return false
			}
			isHealthy, earlyReturn := m.keeper.LiquidatePositionIfUnhealthy(ctx, &position, pool, ammPool)
			if !earlyReturn && isHealthy {
				return true
			}
			return false
		})

		// Close stopLossPrice reached positions
		m.keeper.IteratePoolPosIdsStopLossSorted(ctx, pool.AmmPoolId, func(posId types.AddressId) bool {
			position, err := m.keeper.GetLegacyPosition(ctx, posId.Address, posId.Id)
			if err != nil {
				return false
			}
			underStopLossPrice, earlyReturn := m.keeper.ClosePositionIfUnderStopLossPrice(ctx, &position, pool, ammPool)
			if !earlyReturn && underStopLossPrice {
				return true
			}
			return false
		})
	}
	return nil
}
