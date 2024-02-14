package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.InterestRateMin = sdk.NewDecWithPrec(10, 2) // 10%
	params.InterestRateMax = sdk.NewDecWithPrec(50, 2) // 50%
	m.keeper.SetParams(ctx, params)
	return nil
}
