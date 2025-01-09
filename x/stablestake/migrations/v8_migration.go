package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	m.keeper.MoveAllDebt(ctx)
	m.keeper.MoveAllInterest(ctx)
	return nil
}
