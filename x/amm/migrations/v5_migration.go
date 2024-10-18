package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (m Migrator) V5Migration(ctx sdk.Context) error {
	m.keeper.SetParams(ctx, types.DefaultParams())
	return nil
}
