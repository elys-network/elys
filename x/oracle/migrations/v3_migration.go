package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.LifeTimeInBlocks = 4
	m.keeper.SetParams(ctx, params)
	return nil
}
