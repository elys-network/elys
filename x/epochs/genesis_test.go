package epochs_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/elys-network/elys/v4/app"
	"github.com/elys-network/elys/v4/x/epochs"
	"github.com/elys-network/elys/v4/x/epochs/types"
)

func TestEpochsExportGenesis(t *testing.T) {
	app := app.InitElysTestApp(true, t)
	ctx := app.BaseApp.NewContext(true)

	chainStartTime := ctx.BlockTime()
	chainStartHeight := ctx.BlockHeight()

	genesisState := types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "band_epoch",
				Duration:                time.Second * 15,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: 0,
				EpochCountingStarted:    false,
			},
			{
				Identifier:              types.WeekEpochID,
				StartTime:               time.Time{},
				Duration:                time.Hour * 24 * 7,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: 0,
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    false,
			},
			{
				Identifier:              types.DayEpochID,
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: 0,
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    false,
			},
			{
				Identifier:              types.TenDaysEpochID,
				Duration:                time.Second * 864000,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: 0,
				EpochCountingStarted:    false,
			},
			{
				Identifier:              types.FiveMinutesEpochID,
				Duration:                time.Second * 300,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: 0,
				EpochCountingStarted:    false,
			},
		},
	}

	epochs.InitGenesis(ctx, *app.EpochsKeeper, genesisState)

	genesis := epochs.ExportGenesis(ctx, *app.EpochsKeeper)
	require.Len(t, genesis.Epochs, 5)

	require.Equal(t, genesis.Epochs[0].Identifier, "band_epoch")
	require.Equal(t, genesis.Epochs[1].Identifier, types.DayEpochID)
	require.Equal(t, genesis.Epochs[1].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[1].Duration, time.Hour*24)
	require.Equal(t, genesis.Epochs[1].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[1].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[1].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[1].EpochCountingStarted, false)
	require.Equal(t, genesis.Epochs[2].Identifier, types.FiveMinutesEpochID)
	require.Equal(t, genesis.Epochs[2].StartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[2].Duration, time.Second*300)
	require.Equal(t, genesis.Epochs[2].CurrentEpoch, int64(0))
	require.Equal(t, genesis.Epochs[2].CurrentEpochStartHeight, chainStartHeight)
	require.Equal(t, genesis.Epochs[2].CurrentEpochStartTime, chainStartTime)
	require.Equal(t, genesis.Epochs[2].EpochCountingStarted, false)
}

func TestEpochsInitGenesis(t *testing.T) {
	app := app.InitElysTestApp(true, t)

	ctx := app.BaseApp.NewContext(true)

	// On init genesis, default epochs information is set
	// To check init genesis again, should make it fresh status
	epochInfos := app.EpochsKeeper.AllEpochInfos(ctx)
	for _, epochInfo := range epochInfos {
		app.EpochsKeeper.DeleteEpochInfo(ctx, epochInfo.Identifier)
	}

	now := time.Now().UTC()
	ctx = ctx.WithBlockHeight(1)
	ctx = ctx.WithBlockTime(now)

	// test genesisState validation
	genesisState := types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}
	require.EqualError(t, genesisState.Validate(), "duplicated epoch entry monthly")

	genesisState = types.GenesisState{
		Epochs: []types.EpochInfo{
			{
				Identifier:              "monthly",
				StartTime:               time.Time{},
				Duration:                time.Hour * 24,
				CurrentEpoch:            0,
				CurrentEpochStartHeight: ctx.BlockHeight(),
				CurrentEpochStartTime:   time.Time{},
				EpochCountingStarted:    true,
			},
		},
	}

	epochs.InitGenesis(ctx, *app.EpochsKeeper, genesisState)
	epochInfo, found := app.EpochsKeeper.GetEpochInfo(ctx, "monthly")
	require.True(t, found)
	require.Equal(t, epochInfo.Identifier, "monthly")
	require.Equal(t, epochInfo.StartTime.UTC().String(), now.UTC().String())
	require.Equal(t, epochInfo.Duration, time.Hour*24)
	require.Equal(t, epochInfo.CurrentEpoch, int64(0))
	require.Equal(t, epochInfo.CurrentEpochStartHeight, ctx.BlockHeight())
	require.Equal(t, epochInfo.CurrentEpochStartTime.UTC().String(), time.Time{}.String())
	require.Equal(t, epochInfo.EpochCountingStarted, true)
}
