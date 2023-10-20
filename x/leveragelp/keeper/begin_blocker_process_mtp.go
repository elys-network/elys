package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/leveragelp/types"
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

	_ = k.SetMTP(ctx, mtp)
	// TODO: missing function
	// repayAmount, err := k.ForceCloseLong(ctx, *mtp, pool, false, true)

	if err == nil {
		// TODO: missing function
		// Emit event if position was closed
		// k.EmitForceClose(ctx, mtp, repayAmount, "")
	} else if err != types.ErrMTPUnhealthy {
		ctx.Logger().Error(errors.Wrap(err, "error executing force close").Error())
	}

}
