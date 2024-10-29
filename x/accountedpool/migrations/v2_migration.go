package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	allAccountedPool := m.keeper.GetAllAccountedPool(ctx)
	for _, accountedPool := range allAccountedPool {
		m.keeper.RemoveAccountedPool(ctx, accountedPool.PoolId)
	}
	return nil
}
