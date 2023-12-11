package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/incentive/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	m.keeper.SetParams(ctx, types.DefaultParams())
	return nil
}
