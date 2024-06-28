package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {
	positions := m.keeper.GetAllPositions(ctx)
	for _, position := range positions {
		m.keeper.DestroyPosition(ctx, position.Address, position.Id)
	}
	return nil
}
