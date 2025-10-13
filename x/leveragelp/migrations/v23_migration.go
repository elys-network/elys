package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V23Migration(ctx sdk.Context) error {
	m.keeper.MigratePositionsToNewKeys(ctx)
	m.keeper.DeleteLegacyFallbackOffset(ctx)

	for _, pool := range m.keeper.GetAllPools(ctx) {
		pool.AdlTriggerRatio = math.LegacyMustNewDecFromStr("0.37")
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
