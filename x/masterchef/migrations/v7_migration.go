package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (m Migrator) V7Migration(ctx sdk.Context) error {
	params := m.keeper.GetParams(ctx)
	params.TakerManager = "elys1sznvjrfpwfatwdd5yxzh5m6fl0kf4zpdkxv5wa"
	m.keeper.SetParams(ctx, params)
	return nil
}
