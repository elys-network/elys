package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V16Migration(ctx sdk.Context) error {
	allPools := m.keeper.GetAllPools(ctx)

	for _, pool := range allPools {
		hooks := m.keeper.GetHooks()

		if hooks != nil {
			ammPool, err := m.keeper.GetAmmPool(ctx, pool.AmmPoolId)
			if err != nil {
				return err
			}
			err = hooks.AfterEnablingPool(ctx, ammPool)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
