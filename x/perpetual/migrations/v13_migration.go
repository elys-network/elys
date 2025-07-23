package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V13Migration(ctx sdk.Context) error {
	//m.keeper.DeleteAllToPay(ctx)
	return nil
}
