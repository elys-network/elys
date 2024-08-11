package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	params := m.keeper.GetLegacyParams(ctx)
	m.keeper.SetParams(ctx, params)
	m.keeper.DeleteLegacyParams(ctx)

	commitments := m.keeper.GetAllLegacyCommitments(ctx)
	for _, c := range commitments {
		// SetCommitments will increase the counter for params.NumberOfCommitments as in HasCommitments c won't be found
		m.keeper.SetCommitments(ctx, *c)
		// DeleteLegacyCommitments check using HasLegacyCommitments, so it will reduce the counter for params.NumberOfCommitments
		m.keeper.DeleteLegacyCommitments(ctx, c.Creator)
	}
	return nil
}
