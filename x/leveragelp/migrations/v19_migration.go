package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V19Migration(ctx sdk.Context) error {
	m.keeper.V19Migration(ctx)
	return nil
}
