package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	//count := m.keeper.GetPositionCount(ctx)
	//m.keeper.SetOpenPositionCount(ctx, count)
	return nil
}
