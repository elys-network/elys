package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V16Migration(ctx sdk.Context) error {
	legacyPools := m.keeper.GetAllLegacyPools(ctx)
	params := m.keeper.GetParams(ctx)

	for _, legacyPool := range legacyPools {
		newPool := types.Pool{
			AmmPoolId:                            legacyPool.AmmPoolId,
			Health:                               legacyPool.Health,
			BorrowInterestRate:                   legacyPool.BorrowInterestRate,
			PoolAssetsLong:                       legacyPool.PoolAssetsLong,
			PoolAssetsShort:                      legacyPool.PoolAssetsShort,
			LastHeightBorrowInterestRateComputed: legacyPool.LastHeightBorrowInterestRateComputed,
			FundingRate:                          legacyPool.FundingRate,
			FeesCollected:                        legacyPool.FeesCollected,
			LeverageMax:                          params.LeverageMax,
		}

		m.keeper.SetPool(ctx, newPool)
	}
	return nil
}
