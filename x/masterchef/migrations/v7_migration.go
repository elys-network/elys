package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.TakerManager = "elys14hj2tavq8fpesdwxxcu44rty3hh90vhujrvcmstl4zr3txmfvw9s29sj82"
	m.keeper.SetParams(ctx, params)
	return nil
}
