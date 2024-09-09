package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {

	params := m.keeper.GetLegacyParams(ctx)
	ctx.Logger().Info("LEGACY PARAMS: ")
	ctx.Logger().Info(params.String())
	ctx.Logger().Info(params.InterestRate.String())
	m.keeper.SetParams(ctx, params)
	newParams := m.keeper.GetParams(ctx)
	ctx.Logger().Info(newParams.String())
	ctx.Logger().Info(newParams.InterestRate.String())
	// Migrate the interest blocks
	interests := m.keeper.GetAllLegacyInterest(ctx)

	for _, interest := range interests {
		m.keeper.SetInterest(ctx, interest.BlockHeight, interest)
	}

	m.keeper.V6_DebtMigration(ctx)

	return nil
}
