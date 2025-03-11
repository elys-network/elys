package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V19Migration(ctx sdk.Context) error {
	m.keeper.SetAllPositions(ctx)
	m.keeper.SetLeveragedAmount(ctx)
	return nil
}
