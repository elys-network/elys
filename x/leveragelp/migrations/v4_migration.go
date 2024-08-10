package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	positions := m.keeper.GetAllPositions(ctx)
	for _, position := range positions {
		m.keeper.SetPosition(ctx, &position)
	}
	return nil
}
