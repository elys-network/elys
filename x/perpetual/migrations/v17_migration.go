package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V17Migration(ctx sdk.Context) error {
	//legacyPools := m.keeper.GetAllLegacyPools(ctx)
	//params := m.keeper.GetParams(ctx)
	//
	//for _, legacyPool := range legacyPools {
	//	newPool := types.Pool{
	//		AmmPoolId:                            legacyPool.AmmPoolId,
	//		BaseAssetLiabilitiesRatio:            legacyPool.Health,
	//		QuoteAssetLiabilitiesRatio:           legacyPool.Health,
	//		BorrowInterestRate:                   legacyPool.BorrowInterestRate,
	//		PoolAssetsLong:                       legacyPool.PoolAssetsLong,
	//		PoolAssetsShort:                      legacyPool.PoolAssetsShort,
	//		LastHeightBorrowInterestRateComputed: legacyPool.LastHeightBorrowInterestRateComputed,
	//		FundingRate:                          legacyPool.FundingRate,
	//		FeesCollected:                        legacyPool.FeesCollected,
	//		LeverageMax:                          params.LeverageMax,
	//	}
	//
	//	err := m.keeper.UpdatePoolHealth(ctx, &newPool)
	//	if err != nil {
	//		return err
	//	}
	//
	//}
	return nil
}
