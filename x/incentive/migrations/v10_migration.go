package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	m.keeper.SetParams(ctx, params)
	return nil
}
