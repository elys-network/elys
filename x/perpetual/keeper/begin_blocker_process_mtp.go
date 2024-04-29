package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func BeginBlockerProcessMTP(ctx sdk.Context, k Keeper, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, baseCurrency string, baseCurrencyDecimal uint64) error {
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
		return errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit liabilities: %s", mtp.String()))
	}
	// calculate and update take profit borrow rate
	mtp.TakeProfitBorrowRate, err = k.CalcMTPTakeProfitBorrowRate(ctx, mtp)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit borrow rate: %s", mtp.String()))
	}
	h, err := k.UpdateMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String()))
	}
	mtp.MtpHealth = h

	// Handle Borrow Interest if within epoch position
	if err := k.HandleBorrowInterest(ctx, mtp, &pool, ammPool); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error handling borrow interest payment: %s", mtp.CollateralAsset))
	}
	if err := k.HandleFundingFeeCollection(ctx, mtp, &pool, ammPool, baseCurrency); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error handling funding fee collection: %s", mtp.CollateralAsset))
	}

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	// check MTP health against threshold
	safetyFactor := k.GetSafetyFactor(ctx)
	var mustForceClose bool = false

	if mtp.MtpHealth.LTE(safetyFactor) {
		// flag position as must force close
		mustForceClose = true
	} else {
		ctx.Logger().Debug(errors.Wrap(types.ErrMTPHealthy, "skipping executing force close because mtp is healthy").Error())
	}

	entry, found := k.assetProfileKeeper.GetEntryByDenom(ctx, mtp.CustodyAsset)
	if !found {
		ctx.Logger().Error(errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", mtp.CustodyAsset).Error())
	}
	custodyAssetDecimals := entry.Decimals

	oneToken := math.NewIntFromBigInt(math.LegacyNewDec(10).Power(uint64(custodyAssetDecimals)).TruncateInt().BigInt())

	assetPrice, err := k.EstimateSwap(ctx, sdk.NewCoin(mtp.CustodyAsset, oneToken), baseCurrency, ammPool)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error estimating swap: %s", mtp.CustodyAsset))
	}

	// divide assetPrice by 10^baseCurrencyDecimal to get the actual price in decimal
	assetPriceDec := math.LegacyNewDecFromBigInt(assetPrice.BigInt()).Quo(math.LegacyNewDec(10).Power(uint64(baseCurrencyDecimal)))

	if types.ReachedTakeProfitPrice(mtp, assetPriceDec) {
		// flag position as must force close
		mustForceClose = true
	} else {
		ctx.Logger().Debug(fmt.Sprintf("skipping force close on position %s because take profit price %s <> %s", mtp.String(), mtp.TakeProfitPrice.String(), assetPrice.String()))
	}

	// if flag is false, then skip force close
	if !mustForceClose {
		return nil
	}

	var repayAmount math.Int
	switch mtp.Position {
	case types.Position_LONG:
		repayAmount, err = k.ForceCloseLong(ctx, mtp, &pool, true, baseCurrency)
	case types.Position_SHORT:
		repayAmount, err = k.ForceCloseShort(ctx, mtp, &pool, true, baseCurrency)
	default:
		return errors.Wrap(types.ErrInvalidPosition, fmt.Sprintf("invalid position type: %s", mtp.Position))
	}

	if err == nil {
		// Emit event if position was closed
		k.EmitForceClose(ctx, mtp, repayAmount, "")
	} else {
		return errors.Wrap(err, "error executing force close")
	}

	return nil
}
