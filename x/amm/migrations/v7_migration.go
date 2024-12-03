package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {

	ctx.Logger().Info("Bank Module Migration Finished")

	return nil
}
