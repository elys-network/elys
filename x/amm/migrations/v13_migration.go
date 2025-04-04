package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V13Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)

	// allowedUpfrontSwapMakers is initialised as empty array

	m.keeper.SetParams(ctx, legacyParams)
	return nil
}
