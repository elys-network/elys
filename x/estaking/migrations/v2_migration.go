package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
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

	legacyParams := m.keeper.GetLegacyParams(ctx)
	m.keeper.SetParams(ctx, legacyParams)
	m.keeper.DeleteLegacyParams(ctx)

	return nil
}
