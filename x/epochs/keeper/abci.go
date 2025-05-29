package keeper

import (
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/telemetry"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/elys-network/elys/v6/x/epochs/types"
)

// BeginBlocker of epochs module
func (k Keeper) BeginBlocker(ctx sdk.Context) {
	defer telemetry.ModuleMeasureSince(types.ModuleName, time.Now(), telemetry.MetricKeyBeginBlocker)

	logger := k.Logger(ctx)

	k.IterateEpochInfo(ctx, func(_ int64, epochInfo types.EpochInfo) (stop bool) {
		// Has it not started, and is the block time > initial epoch start time
		shouldInitialEpochStart := !epochInfo.EpochCountingStarted && !epochInfo.StartTime.After(ctx.BlockTime())

		epochEndTime := epochInfo.CurrentEpochStartTime.Add(epochInfo.Duration)
		shouldEpochEnd := ctx.BlockTime().After(epochEndTime) && !shouldInitialEpochStart && !epochInfo.StartTime.After(ctx.BlockTime())

		epochInfo.CurrentEpochStartHeight = ctx.BlockHeight()

		switch {
		case shouldInitialEpochStart:
			epochInfo.StartInitialEpoch()

			logger.Debug("starting epoch", "identifier", epochInfo.Identifier)
		case shouldEpochEnd:
			epochInfo.EndEpoch()

			logger.Debug("ending epoch", "identifier", epochInfo.Identifier)

			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeEpochEnd,
					sdk.NewAttribute(types.AttributeEpochNumber, strconv.FormatInt(epochInfo.CurrentEpoch, 10)),
				),
			)
			err := k.RunHooksAfterEpochEnd(ctx, epochInfo.Identifier, epochInfo.CurrentEpoch)
			if err != nil {
				panic(err)
			}
		default:
			// continue
			return false
		}

		k.SetEpochInfo(ctx, epochInfo)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeEpochStart,
				sdk.NewAttribute(types.AttributeEpochNumber, strconv.FormatInt(epochInfo.CurrentEpoch, 10)),
				sdk.NewAttribute(types.AttributeEpochStartTime, strconv.FormatInt(epochInfo.CurrentEpochStartTime.Unix(), 10)),
			),
		)

		err := k.RunHooksBeforeEpochStart(ctx, epochInfo.Identifier, epochInfo.CurrentEpoch)
		if err != nil {
			panic(err)
		}

		return false
	})
}
