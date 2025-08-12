package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V23Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	for _, pool := range m.keeper.GetAllPools(ctx) {
		pool.MtpSafetyFactor = params.SafetyFactor
		m.keeper.SetPool(ctx, pool)
	}
	return nil
}
