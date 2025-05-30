package migrations

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/epochs/types"
)

func (m Migrator) V3Migration(ctx sdk.Context) error {
	// twelve hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twelvehours",
		Duration:                time.Hour * 12,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   ctx.BlockTime(),
		StartTime:               ctx.BlockTime(),
	})

	// six hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "sixhours",
		Duration:                time.Hour * 6,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   ctx.BlockTime(),
		StartTime:               ctx.BlockTime(),
	})

	// four hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "fourhours",
		Duration:                time.Hour * 4,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   ctx.BlockTime(),
		StartTime:               ctx.BlockTime(),
	})

	// two hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twohours",
		Duration:                time.Hour * 2,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   ctx.BlockTime(),
		StartTime:               ctx.BlockTime(),
	})

	// half hour
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "halfhour",
		Duration:                time.Minute * 30,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
		CurrentEpochStartTime:   ctx.BlockTime(),
		StartTime:               ctx.BlockTime(),
	})

	return nil
}
