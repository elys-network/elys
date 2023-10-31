package keeper

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		currentHeight := ctx.BlockHeight()
		_ = currentHeight
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
			if err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
				continue
			}
			if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
				mtps, _, _ := k.GetMTPsForPool(ctx, pool.AmmPoolId, nil)
				for _, mtp := range mtps {
					k.LiquidatePositionIfUnhealthy(ctx, mtp, pool, ammPool)
				}
			}
			k.SetPool(ctx, pool)
		}
	}

}

func (k Keeper) LiquidatePositionIfUnhealthy(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.GetMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String())).Error())
		return
	}
	mtp.MtpHealth = h
	k.SetMTP(ctx, mtp)

	params := k.GetParams(ctx)
	if mtp.MtpHealth.GT(params.SafetyFactor) {
		return
	}

	repayAmount, err := k.ForceCloseLong(ctx, *mtp, pool)
	if err == nil {
		ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventClose,
			sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
			sdk.NewAttribute("address", mtp.Address),
			sdk.NewAttribute("collateral", mtp.Collateral.String()),
			sdk.NewAttribute("repay_amount", repayAmount.String()),
			sdk.NewAttribute("leverage", mtp.Leverage.String()),
			sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
			sdk.NewAttribute("health", mtp.MtpHealth.String()),
		))
	} else {
		ctx.Logger().Error(errors.Wrap(err, "error executing force close").Error())
	}
}
