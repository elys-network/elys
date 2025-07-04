package keeper

import (
	"errors"
	"fmt"

	sdkerrors "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) CheckAndLiquidatePosition(ctx sdk.Context, mtp *types.MTP, pool types.Pool, ammPool *ammtypes.Pool, closer string) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("function RouteExactAmountOut failed due to internal reason: %v", r)
			ctx.Logger().Error(err.Error())
		}
	}()
	totalPerpetualFeesInUsd := types.NewPerpetualFeesWithEmptyCoins()
	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesInUsd, err := k.MTPTriggerChecksAndUpdates(ctx, mtp, &pool, ammPool)
	if err != nil {
		return err
	}
	totalPerpetualFeesInUsd = perpetualFeesInUsd

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return err
	}

	if forceClosed {
		k.EmitForceClose(ctx, "unhealthy", *mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesInUsd)
		return
	}

	if mtp.CheckForStopLoss(tradingAssetPrice) {
		repayAmt, returnAmt, perpetualFeesInUsd, err = k.ForceClose(ctx, mtp, &pool, ammPool)
		if err != nil {
			return sdkerrors.Wrap(err, "error executing force close")
		}
		totalPerpetualFeesInUsd = totalPerpetualFeesInUsd.Add(perpetualFeesInUsd)
		k.EmitForceClose(ctx, "stop_loss", *mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesInUsd)
		return
	}

	if mtp.CheckForTakeProfit(tradingAssetPrice) {
		repayAmt, returnAmt, perpetualFeesInUsd, err = k.ForceClose(ctx, mtp, &pool, ammPool)
		if err != nil {
			return sdkerrors.Wrap(err, "error executing force close")
		}
		totalPerpetualFeesInUsd = totalPerpetualFeesInUsd.Add(perpetualFeesInUsd)
		k.EmitForceClose(ctx, "take_profit", *mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, closer, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesInUsd)
		return
	}
	err = errors.New("position cannot be liquidated")
	return
}
