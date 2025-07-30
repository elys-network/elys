package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPools(ctx)
	// Reset pools
	for _, pool := range pools {
		pool.LeveragedLpAmount = math.NewInt(0)
		m.keeper.SetPool(ctx, pool)
	}

	//m.keeper.MigrateData(ctx)
	return nil
}
