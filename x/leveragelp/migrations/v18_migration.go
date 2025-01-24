package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V18Migration(ctx sdk.Context) error {
	m.keeper.V18MigratonPoolLiabilities(ctx)
	return nil
}
