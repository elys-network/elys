package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {

	m.keeper.DeleteAllPendingSpotOrder(ctx)

	return nil
}
