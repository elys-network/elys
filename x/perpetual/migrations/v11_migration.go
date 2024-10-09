package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	pools := m.keeper.GetAllLegacyPools(ctx)

	ctx.Logger().Info("Migrating positions from legacy to new format")

	for _, pool := range pools {
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
			FeesCollected:                        sdk.Coins{},
		}
		m.keeper.RemoveLegacyPool(ctx, newPool.AmmPoolId)
		m.keeper.SetPool(ctx, newPool)
	}

	return nil
}
