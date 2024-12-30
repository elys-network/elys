package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	allCommitments := m.keeper.GetAllCommitments(ctx)

	var totalCommited sdk.Coins
	for _, commitment := range allCommitments {
		for _, token := range commitment.CommittedTokens {
			totalCommited = totalCommited.Add(sdk.NewCoin(token.Denom, token.Amount))
		}
	}

	params := m.keeper.GetParams(ctx)
	params.TotalCommitted = totalCommited
	m.keeper.SetParams(ctx, params)

	return nil
}
