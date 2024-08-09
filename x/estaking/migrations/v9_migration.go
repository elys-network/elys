package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	allElysStaked := m.keeper.GetAllLegacyElysStaked(ctx)
	for _, elysStaked := range allElysStaked {
		m.keeper.SetElysStaked(ctx, elysStaked)
		m.keeper.DeleteLegacyElysStaked(ctx, elysStaked.Address)
	}

	allElysStakeChange := m.keeper.GetAllLegacyElysStakeChange(ctx)
	for _, elysStakeChange := range allElysStakeChange {
		m.keeper.SetElysStakeChange(ctx, elysStakeChange)
		m.keeper.DeleteLegacyElysStakeChange(ctx, elysStakeChange)
	}

	return nil
}
