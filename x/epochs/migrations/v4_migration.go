package migrations

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/epochs/types"
)

func (m Migrator) V4Migration(ctx sdk.Context) error {
	// added time.Now() in purpose to simulate a consensus breaking scenario to make sure software upgrade test can capture the issue
	// this should result on a consensus failure because both node in localnet will have a different app hash
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twelvehours",
		Duration:                time.Hour * 12,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   time.Now(),
		StartTime:               time.Now(),
	})

	return nil
}
