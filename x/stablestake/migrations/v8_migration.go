package migrations

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (m Migrator) V8Migration(ctx sdk.Context) error {
	legacyParams := m.keeper.GetLegacyParams(ctx)
	params := types.Params{
		DepositDenom:         legacyParams.DepositDenom,
		RedemptionRate:       legacyParams.RedemptionRate,
		EpochLength:          legacyParams.EpochLength,
		InterestRate:         legacyParams.InterestRate,
		InterestRateMax:      legacyParams.InterestRateMax,
		InterestRateMin:      legacyParams.InterestRateMin,
		InterestRateIncrease: legacyParams.InterestRateIncrease,
		InterestRateDecrease: legacyParams.InterestRateDecrease,
		HealthGainFactor:     legacyParams.HealthGainFactor,
		TotalValue:           legacyParams.TotalValue,
		MaxLeverageRatio:     legacyParams.MaxLeverageRatio,
		MaxBorrowRatio:       math.LegacyMustNewDecFromStr("0.8"),
	}
	m.keeper.SetParams(ctx, params)

	return nil
}
