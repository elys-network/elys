package migrations

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/epochs/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	// twelve hours
	m.keeper.SetEpochInfo(ctx, types.EightHourEpochInfo)
	return nil
}
