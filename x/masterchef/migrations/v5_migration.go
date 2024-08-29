package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPoolInfos(ctx)
	for _, pool := range pools {
		pool.EnableEdenRewards = true
		m.keeper.SetPoolInfo(ctx, pool)
	}
	return nil
}
