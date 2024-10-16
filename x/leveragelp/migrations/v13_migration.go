package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V13Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPools(ctx)
	// Reset pools
	for _, pool := range pools {
		pool.LeverageMax = sdk.NewDec(60)
		m.keeper.SetPool(ctx, pool)
	}

	return nil
}