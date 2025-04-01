package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V10Migration(ctx sdk.Context) error {
	if ctx.ChainID() == "elys-mainnet" {
		return nil
	}
	m.keeper.V10Migrate(ctx)
	return nil
}
