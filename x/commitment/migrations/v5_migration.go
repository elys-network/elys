package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	commitments := m.keeper.GetAllLegacyCommitments(ctx)
	for _, c := range commitments {
		// SetCommitments will increase the counter for params.NumberOfCommitments as HasCommitments won't be found
		m.keeper.SetCommitments(ctx, *c)
		// DeleteLegacyCommitments check using HasLegacyCommitments, so it will reduce the counter for params.NumberOfCommitments
		m.keeper.DeleteLegacyCommitments(ctx, c.Creator)
	}
	return nil
}
