package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {

	m.keeper.V6_MTPMigration(ctx)

	allLegacyPools := m.keeper.GetAllLegacyPools(ctx)
	for _, pool := range allLegacyPools {
		m.keeper.SetPool(ctx, pool)
		m.keeper.RemoveLegacyPool(ctx, pool.AmmPoolId)
	}

	m.keeper.V6_MigrateWhitelistedAddress(ctx)

	return nil
}
