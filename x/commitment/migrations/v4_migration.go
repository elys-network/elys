package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/commitment/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	commitments := m.keeper.GetAllCommitments(ctx)
	for _, c := range commitments{
		newCommittedTokens := []*types.CommittedTokens{}
		for _, commitments := range c.CommittedTokens {
			if commitments.Amount.LTE(sdk.ZeroInt()){
				newCommittedTokens = append(newCommittedTokens, commitments)
			}
		}
		c.CommittedTokens = newCommittedTokens
		m.keeper.SetCommitments(ctx, *c)
	}
	return nil
}
