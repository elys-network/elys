package migrations

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/leveragelp/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPools(ctx)
	// Reset pools
	for _, pool := range pools {
		pool.LeveragedLpAmount = math.NewInt(0)
		m.keeper.SetPool(ctx, pool)
	}
	// Traverse positions and update lp amount and health, as there are few positions, haven't optimized it much
	positions := m.keeper.GetAllPositions(ctx)
	for _, position := range positions {
		// Retrieve Pool
		pool, found := m.keeper.GetPool(ctx, position.AmmPoolId)
		if !found {
			return errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
		}
		pool.LeveragedLpAmount = pool.LeveragedLpAmount.Add(position.LeveragedLpAmount)
		pool.Health = m.keeper.CalculatePoolHealth(ctx, &pool)
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
