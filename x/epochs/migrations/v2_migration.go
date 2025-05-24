package migrations

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v5/x/epochs/types"
)

func (m Migrator) V2Migration(ctx sdk.Context) error {
	// twelve hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twelvehours",
		Duration:                time.Hour * 12,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
	})

	// six hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "sixhours",
		Duration:                time.Hour * 6,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
	})

	// four hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "fourhours",
		Duration:                time.Hour * 4,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
	})

	// two hours
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "twohours",
		Duration:                time.Hour * 2,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
	})

	// half hour
	m.keeper.SetEpochInfo(ctx, types.EpochInfo{
		Identifier:              "halfhour",
		Duration:                time.Minute * 30,
		CurrentEpoch:            0,
		CurrentEpochStartHeight: 0,
		EpochCountingStarted:    false,
	})

	return nil
}
