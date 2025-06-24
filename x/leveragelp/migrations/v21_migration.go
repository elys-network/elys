package migrations

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/leveragelp/types"
)

func (m Migrator) V21Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPools(ctx)
	for _, pool := range pools {
		if pool.AmmPoolId == 2 || pool.AmmPoolId == 10 || pool.AmmPoolId == 15 {
			pool.LeveragedLpAmount = math.NewInt(0)
			m.keeper.SetPool(ctx, pool)
		}
	}
	// Traverse positions and update lp amount and health, as there are few positions, haven't optimized it much
	positions := m.keeper.GetAllPositions(ctx)
	for _, position := range positions {
		if position.AmmPoolId == 2 || position.AmmPoolId == 10 || position.AmmPoolId == 15 {
			// Retrieve Pool
			pool, found := m.keeper.GetPool(ctx, position.AmmPoolId)
			if !found {
				return errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
			}
			pool.LeveragedLpAmount = pool.LeveragedLpAmount.Add(position.LeveragedLpAmount)
			pool.Health = m.keeper.CalculatePoolHealth(ctx, &pool).Dec()
			m.keeper.SetPool(ctx, pool)
		}
	}
	return nil
}
