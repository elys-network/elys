package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/margin/types"
)

func BeginBlockerProcessMTP(ctx sdk.Context, k Keeper, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, baseCurrency string) {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()
	var err error
	// update mtp take profit liabilities
	// calculate mtp take profit liablities, delta x_tp_l = delta y_tp_c * current price (take profit liabilities = take profit custody * current price)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, mtp, baseCurrency)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit liabilities: %s", mtp.String())).Error())
		return
	}
	// calculate and update take profit borrow rate
	mtp.TakeProfitBorrowRate, err = k.CalcMTPTakeProfitBorrowRate(ctx, mtp)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit borrow rate: %s", mtp.String())).Error())
		return
	}
	h, err := k.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String())).Error())
		return
	}
	mtp.MtpHealth = h
	// compute borrow interest
	// TODO: missing fields
	for _, custody := range mtp.Custodies {
		custodyAsset := custody.Denom
		// Retrieve AmmPool
		ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error retrieving amm pool: %d", mtp.AmmPoolId)).Error())
			return
		}

		for _, collateral := range mtp.Collaterals {
			collateralAsset := collateral.Denom
			// Handle Borrow Interest if within epoch position
			if err := k.HandleBorrowInterest(ctx, mtp, &pool, ammPool, collateralAsset, custodyAsset); err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error handling borrow interest payment: %s", collateralAsset)).Error())
				return
			}
			if err := k.HandleFundingFeeCollection(ctx, mtp, &pool, ammPool, collateralAsset, custodyAsset); err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error handling funding fee collection: %s", collateralAsset)).Error())
				return
			}
		}
	}

	_ = k.SetMTP(ctx, mtp)

	var mustForceClose bool = false
	for _, custody := range mtp.Custodies {
		assetPrice, err := k.EstimateSwap(ctx, sdk.NewCoin(custody.Denom, sdk.OneInt()), baseCurrency, ammPool)
		if err != nil {
			ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error estimating swap: %s", custody.Denom)).Error())
			continue
		}
		if mtp.TakeProfitPrice.GT(sdk.NewDecFromInt(assetPrice)) {
			// flag position as must force close
			mustForceClose = true
			break
		}
		ctx.Logger().Error(fmt.Sprintf("skipping force close on position %s because take profit price %s is less than asset price %s", mtp.String(), mtp.TakeProfitPrice.String(), sdk.NewDecFromInt(assetPrice).String()))
	}

	// check MTP health against threshold
	safetyFactor := k.GetSafetyFactor(ctx)

	if mtp.MtpHealth.GT(safetyFactor) {
		ctx.Logger().Error(errors.Wrap(types.ErrMTPHealthy, "skipping executing force close because mtp is healthy").Error())
	} else {
		// flag position as must force close
		mustForceClose = true
	}

	// if flag is false, then skip force close
	if !mustForceClose {
		return
	}

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
	} else {
		ctx.Logger().Error(errors.Wrap(err, "error executing force close").Error())
	}

}
