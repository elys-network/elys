package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	m.keeper.MoveAllDebt(ctx)
	m.keeper.MoveAllInterest(ctx)

	params := m.keeper.GetParams(ctx)
	pool := types.Pool{
		Id:                   types.UsdcPoolId,
		DepositDenom:         params.LegacyDepositDenom,
		InterestRateDecrease: params.LegacyInterestRateDecrease,
		InterestRateIncrease: params.LegacyInterestRateIncrease,
		HealthGainFactor:     params.LegacyHealthGainFactor,
		MaxLeverageRatio:     params.LegacyMaxLeverageRatio,
		MaxWithdrawRatio:     params.LegacyMaxWithdrawRatio,
		InterestRateMax:      params.LegacyInterestRateMax,
		InterestRateMin:      params.LegacyInterestRateMin,
		InterestRate:         params.LegacyInterestRate,
		TotalValue:           params.TotalValue,
	}

	m.keeper.SetPool(ctx, pool)
	return nil
}
