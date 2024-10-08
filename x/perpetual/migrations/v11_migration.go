package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V11Migration(ctx sdk.Context) error {
	mtps := m.keeper.GetAllMTPs(ctx)
	for _, mtp := range mtps {
		m.keeper.CheckSamePositionAndConsolidate(ctx, &mtp)
	}

	return nil
}