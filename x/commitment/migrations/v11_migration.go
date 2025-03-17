package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	m.keeper.V11_SetVested(ctx)
	return nil
}
