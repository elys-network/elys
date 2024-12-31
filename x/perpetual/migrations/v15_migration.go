package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V15Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	newParams := types.Params{
		LeverageMax:                         legacyParams.LeverageMax,
		BorrowInterestRateMax:               legacyParams.BorrowInterestRateMax,
		BorrowInterestRateMin:               legacyParams.BorrowInterestRateMin,
		BorrowInterestRateIncrease:          legacyParams.BorrowInterestRateIncrease,
		BorrowInterestRateDecrease:          legacyParams.BorrowInterestRateDecrease,
		HealthGainFactor:                    legacyParams.HealthGainFactor,
		MaxOpenPositions:                    legacyParams.MaxOpenPositions,
		PoolOpenThreshold:                   legacyParams.PoolOpenThreshold,
		BorrowInterestPaymentFundPercentage: legacyParams.BorrowInterestPaymentFundPercentage,
		BorrowInterestPaymentFundAddress:    legacyParams.BorrowInterestPaymentFundAddress,
		SafetyFactor:                        legacyParams.SafetyFactor,
		BorrowInterestPaymentEnabled:        legacyParams.BorrowInterestPaymentEnabled,
		WhitelistingEnabled:                 legacyParams.WhitelistingEnabled,
		PerpetualSwapFee:                    legacyParams.PerpetualSwapFee,
		MaxLimitOrder:                       legacyParams.MaxLimitOrder,
		FixedFundingRate:                    legacyParams.FixedFundingRate,
		MinimumLongTakeProfitPriceRatio:     legacyParams.MinimumLongTakeProfitPriceRatio,
		MaximumLongTakeProfitPriceRatio:     legacyParams.MaximumLongTakeProfitPriceRatio,
		MaximumShortTakeProfitPriceRatio:    legacyParams.MaximumShortTakeProfitPriceRatio,
		EnableTakeProfitCustodyLiabilities:  legacyParams.EnableTakeProfitCustodyLiabilities,
		WeightBreakingFeeFactor:             legacyParams.WeightBreakingFeeFactor,
		EnabledPools:                        []uint64{},
	}

	err := m.keeper.SetParams(ctx, &newParams)
	if err != nil {
		return err
	}
	return nil
}
