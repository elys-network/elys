package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) HandleInterestPayment(ctx sdk.Context, interestPayment sdk.Int, mtp *types.MTP, pool *types.Pool, ammPool ammtypes.Pool) sdk.Int {
	incrementalInterestPaymentEnabled := k.GetIncrementalInterestPaymentEnabled(ctx)
	// if incremental payment on, pay interest
	if incrementalInterestPaymentEnabled {
		finalInterestPayment, err := k.IncrementalInterestPayment(ctx, interestPayment, mtp, pool, ammPool)
		if err != nil {
			ctx.Logger().Error(sdkerrors.Wrap(err, "error executing incremental interest payment").Error())
		} else {
			return finalInterestPayment
		}
	} else { // else update unpaid mtp interest
		mtp.InterestUnpaidCollateral = interestPayment
	}
	return sdk.ZeroInt()
}
