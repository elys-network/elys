package keeper

import (
	sdkerrors "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) CheckAndLiquidatePosition(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool *ammtypes.Pool, closer string) error {
	defer func() {
		if r := recover(); r != nil {
			if msg, ok := r.(string); ok {
				ctx.Logger().Error(msg)
			}
		}
	}()

	repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.MTPTriggerChecksAndUpdates(ctx, mtp, &pool, ammPool)
	if err != nil {
		return err
	}

	tradingAssetPrice, err := k.GetAssetPrice(ctx, mtp.TradingAsset)
	if err != nil {
		return err
	}

	if forceClosed {
		k.EmitForceClose(ctx, "unhealthy", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
		return nil
	}

	if mtp.CheckForStopLoss(tradingAssetPrice) {
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, &pool, ammPool)
		if err != nil {
			return sdkerrors.Wrap(err, "error executing force close")
		}
		k.EmitForceClose(ctx, "stop_loss", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
		return nil
	}

	if mtp.CheckForTakeProfit(tradingAssetPrice) {
		repayAmt, returnAmt, err = k.ForceClose(ctx, mtp, &pool, ammPool)
		if err != nil {
			return sdkerrors.Wrap(err, "error executing force close")
		}
		k.EmitForceClose(ctx, "take_profit", *mtp, repayAmt, returnAmt, fundingFeeAmt, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice)
		return nil
	}

	return errors.New("position cannot be liquidated")
}
