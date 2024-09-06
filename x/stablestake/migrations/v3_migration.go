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

	debts := m.keeper.GetAllLegacyDebts(ctx)
	for _, debt := range debts {
		m.keeper.SetDebt(ctx, debt)
		m.keeper.DeleteLegacyDebt(ctx, debt)
	}

	return nil
}
