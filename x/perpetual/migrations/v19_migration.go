package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (m Migrator) V19Migration(ctx sdk.Context) error {

	legacyParams := m.keeper.GetLegacyParams(ctx)

	params := types.Params{
		LeverageMax:                         legacyParams.LeverageMax,
		BorrowInterestRateMax:               legacyParams.BorrowInterestRateMax,
		BorrowInterestRateMin:               legacyParams.BorrowInterestRateMin,
		BorrowInterestRateIncrease:          legacyParams.BorrowInterestRateIncrease,
		BorrowInterestRateDecrease:          legacyParams.BorrowInterestRateDecrease,
		HealthGainFactor:                    legacyParams.HealthGainFactor,
		MaxOpenPositions:                    legacyParams.MaxOpenPositions,
		PoolMaxLiabilitiesThreshold:         math.LegacyMustNewDecFromStr("0.3"),
		BorrowInterestPaymentFundPercentage: legacyParams.BorrowInterestPaymentFundPercentage,
		SafetyFactor:                        legacyParams.SafetyFactor,
		BorrowInterestPaymentEnabled:        legacyParams.BorrowInterestPaymentEnabled,
		WhitelistingEnabled:                 true,
		PerpetualSwapFee:                    legacyParams.PerpetualSwapFee,
		MaxLimitOrder:                       legacyParams.MaxLimitOrder,
		FixedFundingRate:                    legacyParams.FixedFundingRate,
		MinimumLongTakeProfitPriceRatio:     legacyParams.MinimumLongTakeProfitPriceRatio,
		MaximumLongTakeProfitPriceRatio:     legacyParams.MaximumLongTakeProfitPriceRatio,
		MaximumShortTakeProfitPriceRatio:    legacyParams.MaximumShortTakeProfitPriceRatio,
		WeightBreakingFeeFactor:             legacyParams.WeightBreakingFeeFactor,
		EnabledPools:                        legacyParams.EnabledPools,
		MinimumNotionalValue:                math.LegacyNewDec(10),
		LongMinimumLiabilityAmount:          math.NewInt(1_000_000),
	}

	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	allLegacyPools := m.keeper.GetAllLegacyPools(ctx)
	for _, legacyPool := range allLegacyPools {
		var poolAssetLong []types.PoolAsset
		for _, legacyPoolAsset := range legacyPool.PoolAssetsLong {
			poolAssetLong = append(poolAssetLong, types.PoolAsset{
				AssetDenom:  legacyPoolAsset.AssetDenom,
				Liabilities: legacyPoolAsset.Liabilities,
				Custody:     legacyPoolAsset.Custody,
				Collateral:  legacyPoolAsset.Collateral,
			})
		}

		var poolAssetShort []types.PoolAsset
		for _, legacyPoolAsset := range legacyPool.PoolAssetsShort {
			poolAssetShort = append(poolAssetShort, types.PoolAsset{
				AssetDenom:  legacyPoolAsset.AssetDenom,
				Liabilities: legacyPoolAsset.Liabilities,
				Custody:     legacyPoolAsset.Custody,
				Collateral:  legacyPoolAsset.Collateral,
			})
		}
		pool := types.Pool{
			AmmPoolId:                            legacyPool.AmmPoolId,
			BaseAssetLiabilitiesRatio:            legacyPool.BaseAssetLiabilitiesRatio,
			QuoteAssetLiabilitiesRatio:           legacyPool.QuoteAssetLiabilitiesRatio,
			BorrowInterestRate:                   legacyPool.BorrowInterestRate,
			PoolAssetsLong:                       poolAssetLong,
			PoolAssetsShort:                      poolAssetShort,
			LastHeightBorrowInterestRateComputed: legacyPool.LastHeightBorrowInterestRateComputed,
			FundingRate:                          legacyPool.FundingRate,
			FeesCollected:                        legacyPool.FeesCollected,
			LeverageMax:                          legacyPool.LeverageMax,
		}
		m.keeper.SetPool(ctx, pool)
	}
	err = m.keeper.ResetStore(ctx)
	if err != nil {
		return err
	}
	return nil
}
