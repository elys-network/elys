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
	// interestPayment := k.CalcMTPInterestLiabilities(ctx, *mtp, pool.InterestRate, 0, 0, ammPool, mtp.CollateralAssets[0])
	// finalInterestPayment := k.HandleInterestPayment(ctx, mtp.CollateralAssets[0],mtp.CustodyAssets[0], interestPayment, mtp, &pool, ammPool)
	// nativeAsset := types.GetSettlementAsset()
	// if types.StringCompare(mtp.CollateralAsset, nativeAsset) { // custody is external, payment is custody
	// 	pool.BlockInterestExternal = pool.BlockInterestExternal.Add(finalInterestPayment)
	// } else { // custody is native, payment is custody
	// 	pool.BlockInterestNative = pool.BlockInterestNative.Add(finalInterestPayment)
	// }

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
