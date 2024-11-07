package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (m Migrator) V12Migration(ctx sdk.Context) error {
	params := types.DefaultParams()
	err := m.keeper.SetParams(ctx, &params)
	if err != nil {
		return err
	}
	return nil
}
