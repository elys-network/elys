package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmenttypes "github.com/elys-network/elys/x/commitment/types"
	epochtypes "github.com/elys-network/elys/x/epochs/types"
	"github.com/elys-network/elys/x/launchpad/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) SetEpochInfo(ctx sdk.Context) {
	_, found := k.epochKeeper.GetEpochInfo(ctx, "day")
	if found {
		return
	}

	k.epochKeeper.SetEpochInfo(ctx, epochtypes.EpochInfo{
		Identifier:              "day",
		StartTime:               ctx.BlockTime(),
		Duration:                time.Hour * 24,
		CurrentEpoch:            0,
		CurrentEpochStartTime:   time.Time{},
		EpochCountingStarted:    false,
		CurrentEpochStartHeight: 0,
	})
}

func (k Keeper) SetElysVestingInfo(ctx sdk.Context) {
	params := k.GetParams(ctx)
	commParams := k.commitmentKeeper.GetParams(ctx)
	for _, vi := range commParams.VestingInfos {
		if vi.BaseDenom == ptypes.Elys {
			return
		}
	}

	// TODO: should consider params.BonusInfo.LockDuration

	commParams.VestingInfos = append(commParams.VestingInfos, &commitmenttypes.VestingInfo{
		BaseDenom:       ptypes.Elys,
		EpochIdentifier: "day",
		NumEpochs:       int64(params.BonusInfo.VestingDuration / 86400),
		NumMaxVestings:  10000,
		VestNowFactor:   sdk.ZeroInt(),
		VestingDenom:    ptypes.Elys,
	})
	k.commitmentKeeper.SetParams(ctx, commParams)
}

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// start vesting after return end time gone
	params := k.GetParams(ctx)
	returnEndTime := params.LaunchpadStarttime + params.LaunchpadDuration + params.ReturnDuration
	if returnEndTime > uint64(ctx.BlockTime().Unix()) {
		return
	}

	k.SetEpochInfo(ctx)
	k.SetElysVestingInfo(ctx)

	allOrders := k.GetAllOrders(ctx)
	for _, order := range allOrders {
		if order.VestingStarted {
			continue
		}
		cacheCtx, write := ctx.CacheContext()
		coins := sdk.Coins{sdk.NewCoin(ptypes.Elys, order.BonusAmount)}
		if err := k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sdk.MustAccAddressFromBech32(order.OrderMaker), coins); err != nil {
			ctx.Logger().Error("failed to send bonus coins to order maker", "err", err)
			continue
		}

		if err := k.commitmentKeeper.DepositLiquidTokensClaimed(cacheCtx, ptypes.Elys, order.BonusAmount, order.OrderMaker); err != nil {
			ctx.Logger().Error("failed to deposit bonus coins to claimed", "err", err)
			continue
		}

		if err := k.commitmentKeeper.ProcessTokenVesting(cacheCtx, ptypes.Elys, order.BonusAmount, order.OrderMaker); err != nil {
			ctx.Logger().Error("failed to process vesting for elys token", "err", err)
			continue
		}
		write()
		order.VestingStarted = true
		k.SetOrder(ctx, order)
	}
}
