package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V20Migration(ctx sdk.Context) error {
	m.keeper.SetLeveragedAmount(ctx)
	return nil
}
