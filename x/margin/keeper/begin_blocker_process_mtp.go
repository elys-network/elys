package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func BeginBlockerProcessMTP(ctx sdk.Context, k Keeper, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	h, err := k.UpdateMTPHealth(ctx, *mtp, ammPool)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String())).Error())
		return
	}
	mtp.MtpHealth = h
	// compute interest
	// TODO: missing fields
	for _, custody := range mtp.Custodies {
		custodyAsset := custody.Denom
		// Retrieve AmmPool
		ammPool, err := k.CloseLongChecker.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error retrieving amm pool: %d", mtp.AmmPoolId)).Error())
			return
		}

		for _, collateral := range mtp.Collaterals {
			collateralAsset := collateral.Denom
			// Handle Interest if within epoch position
			if err := k.CloseLongChecker.HandleInterest(ctx, mtp, &pool, ammPool, collateralAsset, custodyAsset); err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error handling interest payment: %s", collateralAsset)).Error())
				return
			}
		}
	}

	_ = k.SetMTP(ctx, mtp)

	var repayAmount sdk.Int
	switch mtp.Position {
	case types.Position_LONG:
		repayAmount, err = k.ForceCloseLong(ctx, mtp, &pool, true)
	case types.Position_SHORT:
		repayAmount, err = k.ForceCloseShort(ctx, mtp, &pool, true)
	default:
		ctx.Logger().Error(errors.Wrap(types.ErrInvalidPosition, fmt.Sprintf("invalid position type: %s", mtp.Position)).Error())
	}

	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, repayAmount, "")
	} else if err != types.ErrMTPUnhealthy {
		ctx.Logger().Error(errors.Wrap(err, "error executing force close").Error())
	}

}
