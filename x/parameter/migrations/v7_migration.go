package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	// reset params
	params := m.keeper.GetParams(ctx)
	params.EnableTakerFeeSwap = true
	params.TakerFeeCollectionInterval = 1000 // blocks

	m.keeper.SetParams(ctx, params)

	return nil
}
