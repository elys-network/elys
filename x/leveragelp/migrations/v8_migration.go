package migrations

import (
	"fmt"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/leveragelp/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {

	// Traverse positions and update lp amount and health
	// Update data structure
	positions := m.keeper.GetAllPositions(ctx)
	pools := m.keeper.GetAllPools(ctx)
	for _, pool := range pools {
		m.keeper.DeletePoolPosIdsLiquidationSorted(ctx, pool.AmmPoolId)
		m.keeper.DeletePoolPosIdsStopLossSorted(ctx, pool.AmmPoolId)
	}
	m.keeper.DeleteCorruptedKeys(ctx)

	openCount := uint64(0)
	for _, position := range positions {
		pool, found := m.keeper.GetPool(ctx, position.AmmPoolId)
		if !found {
			return errorsmod.Wrap(types.ErrInvalidBorrowingAsset, "invalid pool id")
		}
		ammPool, err := m.keeper.GetAmmPool(ctx, pool.AmmPoolId)
		if err != nil {
			ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
			continue
		}
		isHealthy, _ := m.keeper.LiquidatePositionIfUnhealthy(ctx, &position, pool, ammPool)
		if isHealthy {
			openCount++
		}
	}

	m.keeper.SetOpenPositionCount(ctx, openCount)

	// reset params
	legacy := m.keeper.GetLegacyParams(ctx)
	params := types.Params{
		LeverageMax:         legacy.LeverageMax,
		EpochLength:         legacy.EpochLength,
		MaxOpenPositions:    legacy.MaxOpenPositions,
		PoolOpenThreshold:   legacy.PoolOpenThreshold,
		SafetyFactor:        legacy.SafetyFactor,
		WhitelistingEnabled: legacy.WhitelistingEnabled,
		FallbackEnabled:     false,
		NumberPerBlock:      100,
	}
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	// keys migrations after deleting corrupted keys
	positions = m.keeper.GetAllPositions(ctx)
	for _, position := range positions {
		m.keeper.SetPosition(ctx, &position)
		m.keeper.DeleteLegacyPosition(ctx, position.Address, position.Id)
	}

	return nil
}
