package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	params := types.NewParams()

	params.LeverageMax = legacyParams.LeverageMax
	// params.FundingFeeMinRate = legacyParams.FundingFeeMinRate
	// params.FundingFeeMaxRate = legacyParams.FundingFeeMaxRate
	// params.FundingFeeBaseRate = legacyParams.FundingFeeBaseRate
	params.TakeProfitBorrowInterestRateMin = legacyParams.TakeProfitBorrowInterestRateMin
	params.BorrowInterestRateDecrease = legacyParams.BorrowInterestRateDecrease
	params.BorrowInterestRateIncrease = legacyParams.BorrowInterestRateIncrease
	params.BorrowInterestRateMax = legacyParams.BorrowInterestRateMax
	params.BorrowInterestRateMin = legacyParams.BorrowInterestRateMin
	params.MinBorrowInterestAmount = types.NewParams().MinBorrowInterestAmount
	params.EpochLength = legacyParams.EpochLength
	params.ForceCloseFundAddress = legacyParams.ForceCloseFundAddress
	params.ForceCloseFundPercentage = legacyParams.ForceCloseFundPercentage
	params.HealthGainFactor = legacyParams.HealthGainFactor
	params.IncrementalBorrowInterestPaymentEnabled = legacyParams.IncrementalBorrowInterestPaymentEnabled
	params.IncrementalBorrowInterestPaymentFundAddress = legacyParams.IncrementalBorrowInterestPaymentFundAddress
	params.IncrementalBorrowInterestPaymentFundPercentage = legacyParams.IncrementalBorrowInterestPaymentFundPercentage
	params.InvariantCheckEpoch = legacyParams.InvariantCheckEpoch
	params.LeverageMax = legacyParams.LeverageMax
	params.MaxOpenPositions = legacyParams.MaxOpenPositions
	params.PoolOpenThreshold = legacyParams.PoolOpenThreshold
	params.SafetyFactor = legacyParams.SafetyFactor
	params.WhitelistingEnabled = legacyParams.WhitelistingEnabled
	params.MaxLimitOrder = types.NewParams().MaxLimitOrder

	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}

	return nil
}
