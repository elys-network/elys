package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V14Migration(ctx sdk.Context) error {
	// legacyParams := m.keeper.GetLegacyParams(ctx)
	// newParams := types.Params{
	// 	LeverageMax:                         legacyParams.LeverageMax,
	// 	BorrowInterestRateMax:               legacyParams.BorrowInterestRateMax,
	// 	BorrowInterestRateMin:               legacyParams.BorrowInterestRateMin,
	// 	BorrowInterestRateIncrease:          legacyParams.BorrowInterestRateIncrease,
	// 	BorrowInterestRateDecrease:          legacyParams.BorrowInterestRateDecrease,
	// 	HealthGainFactor:                    legacyParams.HealthGainFactor,
	// 	MaxOpenPositions:                    legacyParams.MaxOpenPositions,
	// 	PoolOpenThreshold:                   legacyParams.PoolOpenThreshold,
	// 	BorrowInterestPaymentFundPercentage: legacyParams.IncrementalBorrowInterestPaymentFundPercentage,
	// 	BorrowInterestPaymentFundAddress:    legacyParams.IncrementalBorrowInterestPaymentFundAddress,
	// 	SafetyFactor:                        legacyParams.SafetyFactor,
	// 	BorrowInterestPaymentEnabled:        legacyParams.IncrementalBorrowInterestPaymentEnabled,
	// 	WhitelistingEnabled:                 legacyParams.WhitelistingEnabled,
	// 	PerpetualSwapFee:                    legacyParams.PerpetualSwapFee,
	// 	MaxLimitOrder:                       legacyParams.MaxLimitOrder,
	// 	FixedFundingRate:                    legacyParams.FixedFundingRate,
	// 	MinimumLongTakeProfitPriceRatio:     legacyParams.MinimumLongTakeProfitPriceRatio,
	// 	MaximumLongTakeProfitPriceRatio:     legacyParams.MaximumLongTakeProfitPriceRatio,
	// 	MaximumShortTakeProfitPriceRatio:    legacyParams.MaximumShortTakeProfitPriceRatio,
	// 	EnableTakeProfitCustodyLiabilities:  legacyParams.EnableTakeProfitCustodyLiabilities,
	// 	WeightBreakingFeeFactor:             legacyParams.WeightBreakingFeeFactor,
	// }

	// err := m.keeper.SetParams(ctx, &newParams)
	// if err != nil {
	// 	return err
	// }
	return nil
}
