package keeper

import (
	"context"
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k msgServer) UpdateTakeProfitPrice(goCtx context.Context, msg *types.MsgUpdateTakeProfitPrice) (*types.MsgUpdateTakeProfitPriceResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	// Load existing mtp
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	mtp, err := k.GetMTP(ctx, msg.PoolId, creator, msg.Id)
	if err != nil {
		return nil, err
	}

	initialCollateralCoin := sdk.NewCoin(mtp.CollateralAsset, mtp.Collateral)

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
		k.EmitForceClose(ctx, "update_take_profit", mtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, msg.Creator, allInterestsPaid, tradingAssetPrice, totalPerpetualFeesCoins, closingPrice, initialCollateralCoin, usdcPrice)
		return &types.MsgUpdateTakeProfitPriceResponse{}, nil
	}

	params := k.GetParams(ctx)
	ratio := msg.Price.Quo(tradingAssetPrice)
	if mtp.Position == types.Position_LONG {
		if ratio.LT(params.MinimumLongTakeProfitPriceRatio) || ratio.GT(params.MaximumLongTakeProfitPriceRatio) {
			return nil, fmt.Errorf("take profit price should be between %s and %s times of current market price for long (current ratio: %s)", params.MinimumLongTakeProfitPriceRatio.String(), params.MaximumLongTakeProfitPriceRatio.String(), ratio.String())
		}
	}
	if mtp.Position == types.Position_SHORT {
		if ratio.GT(params.MaximumShortTakeProfitPriceRatio) {
			return nil, fmt.Errorf("take profit price should be less than %s times of current market price for short (current ratio: %s)", params.MaximumShortTakeProfitPriceRatio.String(), ratio.String())
		}
	}

	mtp.TakeProfitPrice = msg.Price

	err = k.SetMTP(ctx, &mtp)
	if err != nil {
		return nil, err
	}

	k.SetPool(ctx, pool)

	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator)
		if err != nil {
			return nil, err
		}
	}

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpetualFeesCoins)
	interestAmtInUSD := k.amm.CalculateUSDValue(ctx, mtp.CustodyAsset, interestAmt).Dec()

	event := sdk.NewEvent(types.EventUpdateTakeProfitPrice,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("owner", mtp.Address),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
		sdk.NewAttribute("take_profit_price", mtp.TakeProfitPrice.String()),
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

	return &types.MsgUpdateTakeProfitPriceResponse{}, nil
}
