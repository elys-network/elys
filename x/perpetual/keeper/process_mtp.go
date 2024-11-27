package keeper

import (
	"fmt"

	"cosmossdk.io/errors"
	"cosmossdk.io/math"

	sdk "github.com/cosmos/cosmos-sdk/types"

	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckAndLiquidateUnhealthyPosition(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, baseCurrency string) error {
	var err error

	// update mtp take profit liabilities
	// calculate mtp take profit liabilities, delta x_tp_l = delta y_tp_c * current price (take profit liabilities = take profit custody * current price)
	mtp.TakeProfitLiabilities, err = k.CalcMTPTakeProfitLiability(ctx, *mtp)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit liabilities: %s", mtp.String()))
	}
	// calculate and update take profit borrow rate
	err = mtp.UpdateMTPTakeProfitBorrowFactor()
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error calculating mtp take profit borrow rate: %s", mtp.String()))
	}

	k.UpdateMTPBorrowInterestUnpaidLiability(ctx, mtp)
	// Handle Borrow Interest if within epoch position
	if _, err := k.SettleMTPBorrowInterestUnpaidLiability(ctx, mtp, &pool, ammPool); err != nil {
		return errors.Wrap(err, fmt.Sprintf("error handling borrow interest payment: %s", mtp.CollateralAsset))
	}

	err = k.SettleFunding(ctx, mtp, &pool, ammPool)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error handling funding fee: %s", mtp.CollateralAsset))
	}

	h, err := k.GetMTPHealth(ctx, *mtp, ammPool, baseCurrency)
	if err != nil {
		return errors.Wrap(err, fmt.Sprintf("error updating mtp health: %s", mtp.String()))
	}
	mtp.MtpHealth = h

	err = k.SetMTP(ctx, mtp)
	if err != nil {
		return err
	}

	k.SetPool(ctx, pool)

	// check MTP health against threshold
	safetyFactor := k.GetSafetyFactor(ctx)

	if mtp.MtpHealth.LTE(safetyFactor) {
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
			k.EmitForceClose(ctx, types.EventForceCloseUnhealthy, mtp, repayAmount, "")
		} else {
			return errors.Wrap(err, "error executing force close")
		}
	} else {
		ctx.Logger().Debug(errors.Wrap(types.ErrMTPHealthy, "skipping executing force close because mtp is healthy").Error())
	}

	return nil
}

func (k Keeper) CheckAndCloseAtStopLoss(ctx sdk.Context, mtp *types.MTP, pool types.Pool, baseCurrency string) error {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return err
	}

	if mtp.Position == types.Position_LONG {
		underStopLossPrice := !mtp.StopLossPrice.IsNil() && tradingAssetPrice.LTE(mtp.StopLossPrice)
		if !underStopLossPrice {
			return fmt.Errorf("mtp stop loss price is not <=  token price")
		}
	} else {
		underStopLossPrice := !mtp.StopLossPrice.IsNil() && tradingAssetPrice.GTE(mtp.StopLossPrice)
		if !underStopLossPrice {
			return fmt.Errorf("mtp stop loss price is not =>  token price")
		}
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
		k.EmitForceClose(ctx, types.EventForceCloseStopLoss, mtp, repayAmount, "")
	} else {
		return errors.Wrap(err, "error executing force close")
	}

	return nil
}

func (k Keeper) CheckAndCloseAtTakeProfit(ctx sdk.Context, mtp *types.MTP, pool types.Pool, baseCurrency string) error {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return err
	}

	if mtp.Position == types.Position_LONG {
		if !tradingAssetPrice.GTE(mtp.TakeProfitPrice) {
			return fmt.Errorf("mtp take profit price is not >=  token price")
		}
	} else {
		if !tradingAssetPrice.LTE(mtp.TakeProfitPrice) {
			return fmt.Errorf("mtp take profit price is not <=  token price")
		}
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
		k.EmitForceClose(ctx, types.EventForceCloseTakeprofit, mtp, repayAmount, "")
	} else {
		return errors.Wrap(err, "error executing force close")
	}

	return nil
}
