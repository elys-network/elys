package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	m.keeper.V10_SetSupply(ctx)
	return nil
}
