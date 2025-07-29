package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	closedMtp, repayAmount, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, totalPerpetualFeesCoins, closingPrice, initialCollateral, initialCustody, initialLiabilities, err := k.ClosePosition(ctx, msg)
	if err != nil {
		return nil, err
	}

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, closedMtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpetualFeesCoins)
	interestAmtInUSD := k.amm.CalculateUSDValue(ctx, closedMtp.CustodyAsset, interestAmt).Dec()

	netPnLInUSD := k.CalcNetPnLAtClosing(ctx, returnAmt, closedMtp.CustodyAsset, initialCollateral, closingRatio)
	usdcPrice, err := k.GetUSDCPrice(ctx)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventClose,
			sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(closedMtp.Id), 10)),
			sdk.NewAttribute("owner", closedMtp.Address),
			sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(closedMtp.AmmPoolId), 10)),
			sdk.NewAttribute("net_pnl", netPnLInUSD.String()),
			sdk.NewAttribute("collateral_asset", closedMtp.CollateralAsset),
			sdk.NewAttribute("initial_collateral_amount", initialCollateral.Amount.String()),
			sdk.NewAttribute("final_collateral_amount", closedMtp.Collateral.String()),
			sdk.NewAttribute("initial_custody_amount", initialCustody.String()),
			sdk.NewAttribute("final_custody_amount", closedMtp.Custody.String()),
			sdk.NewAttribute("initial_liabilities_amount", initialLiabilities.String()),
			sdk.NewAttribute("final_liabilities_amount", closedMtp.Liabilities.String()),
			sdk.NewAttribute("position", closedMtp.Position.String()),
			sdk.NewAttribute("mtp_health", closedMtp.MtpHealth.String()), // should be there if it's partial close
			sdk.NewAttribute("repay_amount", repayAmount.String()),
			sdk.NewAttribute("return_amount", returnAmt.String()),
			sdk.NewAttribute("funding_fee_amount", fundingFeeAmt.String()),
			sdk.NewAttribute("funding_amount_distributed", fundingAmtDistributed.String()),
			sdk.NewAttribute("interest_amount", interestAmt.String()),
			sdk.NewAttribute("interest_amount_in_usd", interestAmtInUSD.String()),
			sdk.NewAttribute("insurance_amount", insuranceAmt.String()),
			sdk.NewAttribute("funding_fee_paid_custody", closedMtp.FundingFeePaidCustody.String()),
			sdk.NewAttribute("funding_fee_received_custody", closedMtp.FundingFeeReceivedCustody.String()),
			sdk.NewAttribute("borrow_interest_paid_custody", closedMtp.BorrowInterestPaidCustody.String()),
			sdk.NewAttribute("closing_ratio", closingRatio.String()),
			sdk.NewAttribute("closing_price", closingPrice.String()),
			sdk.NewAttribute("open_price", closedMtp.OpenPrice.String()),
			sdk.NewAttribute("usdc_price", usdcPrice.String()),
			sdk.NewAttribute("trading_asset_price", tradingAssetPrice.String()),
			sdk.NewAttribute("all_interests_paid", strconv.FormatBool(allInterestsPaid)), // Funding Fee is fully paid but interest amount is only partially paid then this will be false
			sdk.NewAttribute("force_closed", strconv.FormatBool(forceClosed)),
			sdk.NewAttribute(types.AttributeKeyPerpFee, perpFeesInUsd.String()),
			sdk.NewAttribute(types.AttributeKeySlippage, slippageFeesInUsd.String()),
			sdk.NewAttribute(types.AttributeKeyWeightBreakingFee, weightBreakingFeesInUsd.String()),
			sdk.NewAttribute(types.AttributeTakerFees, takerFeesInUsd.String()),
		))

	return &types.MsgCloseResponse{
		Id:     closedMtp.Id,
		Amount: repayAmount,
	}, nil
}
