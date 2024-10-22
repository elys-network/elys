package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V14Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPools(ctx)

	for _, pool := range pools {

		ammPool, err := m.keeper.GetAmmPool(ctx, pool.AmmPoolId)
		if err != nil {
			return err
		}
		if m.keeper.GetHooks() != nil {
			err = m.keeper.GetHooks().AfterEnablingPool(ctx, ammPool)
			if err != nil {
				return err
			}
		}
	}

	return nil
}
