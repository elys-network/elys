package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	legacyParams.LpLockupDuration = 3600

	m.keeper.SetParams(ctx, legacyParams)
	return nil
}
