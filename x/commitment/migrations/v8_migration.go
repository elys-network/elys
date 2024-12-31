package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	m.keeper.V8_ParamsMigration(ctx)
	return nil
}
