package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	m.keeper.UpdateAllLegacyTimeBasedInflation(ctx)
	return nil
}
