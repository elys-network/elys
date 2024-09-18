package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {

	ctx.Logger().Info("Migrating pool from legacy to new format")

	pools := m.keeper.GetAllLegacyPools(ctx)

	for _, pool := range pools {

		ctx.Logger().Debug("pool", pool)

		newPool := types.Pool{
			AmmPoolId:                            pool.AmmPoolId,
			Health:                               pool.Health,
			Enabled:                              pool.Enabled,
			Closed:                               pool.Closed,
			BorrowInterestRate:                   pool.BorrowInterestRate,
			PoolAssetsLong:                       pool.PoolAssetsLong,
			PoolAssetsShort:                      pool.PoolAssetsShort,
			LastHeightBorrowInterestRateComputed: pool.LastHeightBorrowInterestRateComputed,
			FundingRate:                          pool.FundingRate,
			NetOpenInterest:                      sdk.ZeroInt(),
		}

		m.keeper.RemoveLegacyPool(ctx, pool.AmmPoolId)
		m.keeper.SetPool(ctx, newPool)
	}
	return nil
}
