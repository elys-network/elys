package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	m.keeper.MigrateFromV3UserRewardInfos(ctx)
	return nil
}
