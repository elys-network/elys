package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k msgServer) UpdateStopLoss(goCtx context.Context, msg *types.MsgUpdateStopLoss) (*types.MsgUpdateStopLossResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	return k.Keeper.UpdateStopLoss(ctx, msg)
}

func (k Keeper) UpdateStopLoss(ctx sdk.Context, msg *types.MsgUpdateStopLoss) (*types.MsgUpdateStopLossResponse, error) {
	// Load existing mtp
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, msg.PoolId, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	initialCollateralCoin := sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral)
	initialCustody := mtp.Custody
	initialLiabilities := mtp.Liabilities

	poolId := mtp.AmmPoolId
	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId)
	if err != nil {
		return nil, errorsmod.Wrap(err, "amm pool not found")
	}

	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, err := k.MTPTriggerChecksAndUpdates(ctx, &mtp, &pool, &ammPool)
	if err != nil {
		return nil, err
	}

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, mtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	if forceClosed {
		usdcPrice, err := k.GetUSDCPrice(ctx)
		if err != nil {
			return nil, err
		}
		k.EmitForceClose(ctx, "update_stop_loss", mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, msg.Creator, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesCoins, closingPrice, initialCollateralCoin, initialCustody, initialLiabilities, usdcPrice)
		return &types.MsgUpdateStopLossResponse{}, nil
	}

	if mtp.Position == types.Position_LONG {
		if !msg.Price.IsZero() && msg.Price.GTE(tradingAssetPrice) {
			return nil, fmt.Errorf("stop loss price cannot be greater than equal to tradingAssetPrice for long (Stop loss: %s, asset price: %s)", msg.Price.String(), tradingAssetPrice.String())
		}
	}
	if mtp.Position == types.Position_SHORT {
		if !msg.Price.IsZero() && msg.Price.LTE(tradingAssetPrice) {
			return nil, fmt.Errorf("stop loss price cannot be less than equal to tradingAssetPrice for short (Stop loss: %s, asset price: %s)", msg.Price.String(), tradingAssetPrice.String())
		}
	}

	mtp.StopLossPrice = msg.Price
	err = k.SetMTP(ctx, &mtp)
	if err != nil {
		return nil, err
	}

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpetualFeesCoins)
	interestAmtInUSD := k.amm.CalculateUSDValue(ctx, mtp.CustodyAsset, interestAmt).Dec()

	event := sdk.NewEvent(types.EventUpdateStopLoss,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("owner", mtp.Address),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
		sdk.NewAttribute("stop_loss", mtp.StopLossPrice.String()),
		sdk.NewAttribute("funding_fee_amount", fundingFeeAmt.String()),
		sdk.NewAttribute("interest_amount", interestAmt.String()),
		sdk.NewAttribute("interest_amount_in_usd", interestAmtInUSD.String()),
		sdk.NewAttribute("insurance_amount", insuranceAmt.String()),
		sdk.NewAttribute("funding_fee_paid_custody", mtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", mtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute(types.AttributeKeyPerpFee, perpFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeySlippage, slippageFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeyWeightBreakingFee, weightBreakingFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeTakerFees, takerFeesInUsd.String()),
	)
	ctx.EventManager().EmitEvent(event)

	return &types.MsgUpdateStopLossResponse{}, nil
}
