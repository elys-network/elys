package migrations

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/epochs/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	// twelve hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twelvehours",
		Duration:                time.Hour * 12,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   time.Now(),
		StartTime:               time.Now(),
	})

	// six hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "sixhours",
		Duration:                time.Hour * 6,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   time.Now(),
		StartTime:               time.Now(),
	})

	// four hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "fourhours",
		Duration:                time.Hour * 4,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   time.Now(),
		StartTime:               time.Now(),
	})

	// two hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twohours",
		Duration:                time.Hour * 2,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   time.Now(),
		StartTime:               time.Now(),
	})

	// half hour
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "halfhour",
		Duration:                time.Minute * 30,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   time.Now(),
		StartTime:               time.Now(),
	})

	return nil
}
