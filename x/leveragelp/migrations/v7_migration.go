package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	// Traverse positions and update lp amount and health
	// Update data structure
	//openCount := uint64(0)
	//
	//m.keeper.SetOpenPositionCount(ctx, openCount)
	return nil
}
