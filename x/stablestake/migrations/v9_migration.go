package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (m Migrator) V9Migration(ctx sdk.Context) error {
	allPools := m.keeper.GetAllAmmPools(ctx)
	for _, pool := range allPools {
		pool.TotalLiabilities = sdk.Coins{}
		m.keeper.SetAmmPool(ctx, pool)
	}

	m.keeper.MoveAllDebt(ctx)

	params := m.keeper.GetParams(ctx)
	pool := types.Pool{
		Id:                   types.UsdcPoolId,
		DepositDenom:         m.keeper.GetDepositDenom(ctx),
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
	m.keeper.MoveAllInterest(ctx)

	return nil
}
