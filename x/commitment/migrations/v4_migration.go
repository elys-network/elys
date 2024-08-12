package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
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

	commitments = m.keeper.GetAllCommitments(ctx)
	for _, c := range commitments {
		newCommittedTokens := []*types.CommittedTokens{}
		for _, commitmentToken := range c.CommittedTokens {
			if commitmentToken.Amount.LTE(sdk.ZeroInt()) {
				newCommittedTokens = append(newCommittedTokens, commitmentToken)
			}
		}
		c.CommittedTokens = newCommittedTokens
		m.keeper.SetCommitments(ctx, *c)
	}
	return nil
}
