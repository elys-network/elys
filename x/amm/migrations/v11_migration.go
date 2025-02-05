package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	legacyParams.LpLockupDuration = uint64(3600)

	m.keeper.SetParams(ctx, legacyParams)
	return nil
}
