package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V23Migration(ctx sdk.Context) error {
	m.keeper.MigratePositionsToNewKeys(ctx)
	m.keeper.DeleteLegacyFallbackOffset(ctx)
	return nil
}
