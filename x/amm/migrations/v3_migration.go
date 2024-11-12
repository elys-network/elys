package migrations

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllPool(ctx)
	for _, pool := range pools {
		if pool.PoolParams.UseOracle {
			pool.PoolParams.WeightBreakingFeeMultiplier = sdkmath.LegacyNewDecWithPrec(5, 4) // 0.05%
			m.keeper.SetPool(ctx, pool)
		}
	}
	return nil
}
