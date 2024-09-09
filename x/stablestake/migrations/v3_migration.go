package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	// Migrate the interest blocks
	interests := m.keeper.GetAllLegacyInterest(ctx)

	for _, interest := range interests {
		m.keeper.SetInterest(ctx, interest.BlockHeight, interest)
	}

	m.keeper.V6_DebtMigration(ctx)

	params := m.keeper.GetLegacyParams(ctx)
	m.keeper.SetParams(ctx, params)

	return nil
}
