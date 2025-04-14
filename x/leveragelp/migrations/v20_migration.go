package migrations

import sdk "github.com/cosmos/cosmos-sdk/types"

func (m Migrator) V20Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetParams(ctx)
	legacyParams.StopLossEnabled = false
	return m.keeper.SetParams(ctx, &legacyParams)
}
