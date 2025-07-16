package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.TakerManager = "elys1et8gcr9t58l4xzd28u9kz5f6zj3f2740qtfm6v"
	m.keeper.SetParams(ctx, params)
	return nil
}
