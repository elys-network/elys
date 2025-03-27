package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	m.keeper.V10_SetEdenEdenBSupply(ctx)
	m.keeper.V10_SetEdenVested(ctx)
	return nil
}
