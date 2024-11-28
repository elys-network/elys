package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V6Migration(ctx sdk.Context) error {

	m.keeper.DeleteAllPendingSpotOrder(ctx)
	m.keeper.DeleteAllPendingPerpetualOrder(ctx)

	return nil
}