package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/stablestake/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	oldParams := m.keeper.GetParams(ctx)
	params := types.DefaultParams()
	params.TotalValue = oldParams.TotalValue
	m.keeper.SetParams(ctx, params)
	return nil
}
