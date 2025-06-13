package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	//params := types.DefaultParams()
	//params.BandChannelSource = "channel-0"
	//m.keeper.SetParams(ctx, params)
	//m.keeper.MigrateAllLegacyPrices(ctx)
	return nil
}
