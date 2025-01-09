package keeper

import (
	"errors"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckAndLiquidateUnhealthyPosition(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool *ammtypes.Pool, closer string) error {
	repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, mtp, &pool, ammPool)
	if err != nil {
		return err
	}

	if !forceClosed {
		ctx.Logger().Warn(fmt.Sprintf("skipping executing force close because mtp (id: %d, owner: %s) is healthy. Bot address: %s", mtp.Id, mtp.Address, closer))
	} else {
		tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
		if err != nil {
			return err
		}

		k.EmitForceClose(ctx, "unhealthy", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
	}

	return nil
}

func (k Keeper) CheckAndCloseAtStopLoss(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, closer string) error {
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
			return errors.New("mtp stop loss price is not <=  token price")
		}
	} else {
		underStopLossPrice := !mtp.StopLossPrice.IsNil() && tradingAssetPrice.GTE(mtp.StopLossPrice)
		if !underStopLossPrice {
			return errors.New("mtp stop loss price is not =>  token price")
		}
	}

	repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, mtp, &pool, &ammPool)
	if err != nil {
		return err
	}

	if forceClosed {
		k.EmitForceClose(ctx, "unhealthy", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
		return nil
	} else {
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, &pool, &ammPool)
	}

	if err == nil {
		k.EmitForceClose(ctx, "stop_loss", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
	} else {
		return sdkerrors.Wrap(err, "error executing force close")
	}

	return nil
}

func (k Keeper) CheckAndCloseAtTakeProfit(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool ammtypes.Pool, closer string) error {
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
			return errors.New("mtp take profit price is not >=  token price")
		}
	} else {
		if !tradingAssetPrice.LTE(mtp.TakeProfitPrice) {
			return errors.New("mtp take profit price is not <=  token price")
		}
	}

	repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, mtp, &pool, &ammPool)
	if err != nil {
		return err
	}

	if forceClosed {
		k.EmitForceClose(ctx, "unhealthy", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
		return nil
	} else {
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, &pool, &ammPool)
	}

	if err == nil {
		k.EmitForceClose(ctx, "take_profit", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
	} else {
		return sdkerrors.Wrap(err, "error executing force close")
	}

	return nil
}
