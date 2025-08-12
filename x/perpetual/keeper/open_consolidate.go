package keeper

import (
	"fmt"
	"strconv"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) OpenConsolidate(ctx sdk.Context, existingMtp *types.MTP, newMtp *types.MTP, msg *types.MsgOpen, tradingAsset, baseCurrency string, prevPerpFeesCoins types.PerpetualFees) (*types.MsgOpenResponse, error) {
	poolId := existingMtp.AmmPoolId
	ammPool, err := k.GetAmmPool(ctx, poolId)
	if err != nil {
		return nil, err
	}

	existingMtpCollateralCoin := sdk.NewCoin(existingMtp.CollateralAsset, existingMtp.Collateral)
	initialCustody := existingMtp.Custody
	initialLiabilities := existingMtp.Liabilities

	pool, found := k.GetPool(ctx, poolId)
	if !found {
		return nil, errorsmod.Wrap(types.ErrPoolDoesNotExist, fmt.Sprintf("poolId: %d", poolId))
	}

	repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, perpetualFeesCoins, closingPrice, err := k.MTPTriggerChecksAndUpdates(ctx, existingMtp, &pool, &ammPool)
	if err != nil {
		return nil, err
	}
	totalPerpFeesCoins := perpetualFeesCoins.Add(prevPerpFeesCoins)

	if forceClosed {
		tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, tradingAsset)
		if err != nil {
			return nil, err
		}
		usdcPrice, err := k.GetUSDCPrice(ctx)
		if err != nil {
			return nil, err
		}
		k.EmitForceClose(ctx, "open_consolidate", *existingMtp, repayAmt, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, msg.Creator, allInterestsPaid, tradingAssetPrice, totalPerpFeesCoins, closingPrice, existingMtpCollateralCoin, initialCustody, initialLiabilities, usdcPrice)
		return &types.MsgOpenResponse{
			Id: existingMtp.Id,
		}, nil
	}

	existingMtp, err = k.OpenConsolidateMergeMtp(ctx, existingMtp, newMtp)
	if err != nil {
		return nil, err
	}

	if !newMtp.Liabilities.IsZero() {
		consolidatedOpenPrice := (existingMtp.OpenPrice.MulInt(existingMtp.Custody).Add(newMtp.OpenPrice.MulInt(newMtp.Custody))).QuoInt(existingMtp.Custody.Add(newMtp.Custody))
		existingMtp.OpenPrice = consolidatedOpenPrice
	}

	// overwrite take profit price instead of taking average of both take profit prices
	if msg.TakeProfitPrice.IsPositive() {
		existingMtp.TakeProfitPrice = msg.TakeProfitPrice
	}

	existingMtp.MtpHealth, err = k.GetMTPHealth(ctx, *existingMtp)
	if err != nil {
		return nil, err
	}

	// Check if the MTP is unhealthy
	if existingMtp.MtpHealth.LTE(pool.MtpSafetyFactor) {
		return nil, errorsmod.Wrapf(types.ErrMTPUnhealthy, "(MtpHealth: %s)", existingMtp.MtpHealth.String())
	}

	stopLossPrice := msg.StopLossPrice
	if msg.StopLossPrice.IsNil() || msg.StopLossPrice.IsZero() {
		liquidationPrice, err := k.GetLiquidationPrice(ctx, *existingMtp)
		if err != nil {
			return nil, fmt.Errorf("failed to get liquidation price: %s", err.Error())
		}
		stopLossPrice = liquidationPrice
	}
	existingMtp.StopLossPrice = stopLossPrice

	// Set existing MTP
	if err = k.SetMTP(ctx, existingMtp); err != nil {
		return nil, err
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	if k.hooks != nil {
		// The pool value above was sent in pointer so its updated
		err = k.hooks.AfterPerpetualPositionModified(ctx, ammPool, pool, creator)
		if err != nil {
			return nil, err
		}
	}

	if err = k.CheckLowPoolHealthAndMinimumCustody(ctx, poolId, true); err != nil {
		return nil, err
	}

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpFeesCoins)
	interestAmtInUSD := k.amm.CalculateUSDValue(ctx, existingMtp.CustodyAsset, interestAmt).Dec()

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, existingMtp.TradingAsset)
	if err != nil {
		return nil, err
	}
	usdcPrice, err := k.GetUSDCPrice(ctx)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventOpenConsolidate,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(existingMtp.Id), 10)),
		sdk.NewAttribute("owner", existingMtp.Address),
		sdk.NewAttribute("position", existingMtp.Position.String()),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(existingMtp.AmmPoolId), 10)),
		sdk.NewAttribute("collateral_asset", existingMtp.CollateralAsset),
		sdk.NewAttribute("initial_collateral_amount", existingMtpCollateralCoin.Amount.String()),
		sdk.NewAttribute("final_collateral_amount", existingMtp.Collateral.String()),
		sdk.NewAttribute("initial_custody_amount", initialCustody.String()),
		sdk.NewAttribute("final_custody_amount", existingMtp.Custody.String()),
		sdk.NewAttribute("initial_liabilities_amount", initialLiabilities.String()),
		sdk.NewAttribute("final_liabilities_amount", existingMtp.Liabilities.String()),
		sdk.NewAttribute("liabilities", existingMtp.Liabilities.String()),
		sdk.NewAttribute("new_liabilities", newMtp.Liabilities.String()),
		sdk.NewAttribute("custody", existingMtp.Custody.String()),
		sdk.NewAttribute("new_custody", newMtp.Custody.String()),
		sdk.NewAttribute("interest_amount", interestAmt.String()),
		sdk.NewAttribute("interest_amount_in_usd", interestAmtInUSD.String()),
		sdk.NewAttribute("mtp_health", existingMtp.MtpHealth.String()),
		sdk.NewAttribute("stop_loss_price", existingMtp.StopLossPrice.String()),
		sdk.NewAttribute("take_profit_price", existingMtp.TakeProfitPrice.String()),
		sdk.NewAttribute("take_profit_borrow_factor", existingMtp.TakeProfitBorrowFactor.String()),
		sdk.NewAttribute("funding_fee_paid_custody", existingMtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", existingMtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("open_price", existingMtp.OpenPrice.String()),
		sdk.NewAttribute("trading_asset_price", tradingAssetPrice.String()),
		sdk.NewAttribute("usdc_price", usdcPrice.String()),
		sdk.NewAttribute(types.AttributeKeyPerpFee, perpFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeySlippage, slippageFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeyWeightBreakingFee, weightBreakingFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeTakerFees, takerFeesInUsd.String()),
	))

	return &types.MsgOpenResponse{
		Id: existingMtp.Id,
	}, nil
}
