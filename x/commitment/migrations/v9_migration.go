package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	m.keeper.V9_ParamsMigration(ctx)
	return nil
}
