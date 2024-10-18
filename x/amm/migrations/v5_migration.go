package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllLegacyPool(ctx)
	for _, pool := range pools {
		newPool := types.Pool{
			PoolId:  pool.PoolId,
			Address: pool.Address,
			PoolParams: types.PoolParams{
				SwapFee:                     pool.PoolParams.SwapFee,
				ExitFee:                     pool.PoolParams.ExitFee,
				UseOracle:                   pool.PoolParams.UseOracle,
				WeightBreakingFeeMultiplier: pool.PoolParams.WeightBreakingFeeMultiplier,
				WeightBreakingFeeExponent:   pool.PoolParams.WeightBreakingFeeExponent,
				ExternalLiquidityRatio:      pool.PoolParams.ExternalLiquidityRatio,
				WeightRecoveryFeePortion:    pool.PoolParams.WeightRecoveryFeePortion,
				ThresholdWeightDifference:   pool.PoolParams.ThresholdWeightDifference,
				WeightBreakingFeePortion:    sdk.NewDecWithPrec(50, 2), // 50%
				FeeDenom:                    pool.PoolParams.FeeDenom,
			},
			TotalShares:       pool.TotalShares,
			PoolAssets:        pool.PoolAssets,
			TotalWeight:       pool.TotalWeight,
			RebalanceTreasury: pool.RebalanceTreasury,
		}

		m.keeper.RemovePool(ctx, pool.PoolId)
		m.keeper.SetPool(ctx, newPool)
	}
	return nil
}
