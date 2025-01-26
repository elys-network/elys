package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {

	m.keeper.DeleteAllPendingPerpetualOrder(ctx)
	m.keeper.DeleteAllPendingSpotOrder(ctx)

	return nil
}
