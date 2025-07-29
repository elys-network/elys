package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	params := types.DefaultParams()
	m.keeper.SetParams(ctx, params)
	return nil
}
