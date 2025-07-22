package keeper

import (
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) EmitForceClose(ctx sdk.Context, trigger string, closedMTP types.MTP, repayAmount, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt math.Int, closer string, allInterestsPaid bool, tradingAssetPrice math.LegacyDec, totalPerpetualFeesCoins types.PerpetualFees, closingPrice math.LegacyDec, initialCollateral sdk.Coin, initialCustody, initialLiabilities math.Int, usdcPrice math.LegacyDec) {

	perpFeesInUsd, slippageFeesInUsd, weightBreakingFeesInUsd, takerFeesInUsd := k.GetPerpFeesInUSD(ctx, totalPerpetualFeesCoins)
	netPnLInUSD := k.CalcNetPnLAtClosing(ctx, returnAmt, closedMTP.CustodyAsset, initialCollateral, math.LegacyOneDec())
	interestAmtInUSD := k.amm.CalculateUSDValue(ctx, closedMTP.CustodyAsset, interestAmt).Dec()

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventForceClosed,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(closedMTP.Id), 10)),
		sdk.NewAttribute("owner", closedMTP.Address),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(closedMTP.AmmPoolId), 10)),
		sdk.NewAttribute("position", closedMTP.Position.String()),
		sdk.NewAttribute("net_pnl", netPnLInUSD.String()),
		sdk.NewAttribute("collateral_asset", closedMTP.CollateralAsset),
		sdk.NewAttribute("closer", closer),
		sdk.NewAttribute("initial_collateral_amount", initialCollateral.Amount.String()),
		sdk.NewAttribute("final_collateral_amount", closedMTP.Collateral.String()),
		sdk.NewAttribute("initial_custody_amount", initialCustody.String()),
		sdk.NewAttribute("final_custody_amount", closedMTP.Custody.String()),
		sdk.NewAttribute("initial_liabilities_amount", initialLiabilities.String()),
		sdk.NewAttribute("final_liabilities_amount", closedMTP.Liabilities.String()),
		sdk.NewAttribute("closing_price", closingPrice.String()),
		sdk.NewAttribute("usdc_price", usdcPrice.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("return_amount", returnAmt.String()),
		sdk.NewAttribute("funding_fee_amount", fundingFeeAmt.String()),
		sdk.NewAttribute("funding_amount_distributed", fundingAmtDistributed.String()),
		sdk.NewAttribute("interest_amount", interestAmt.String()),
		sdk.NewAttribute("interest_amount_in_usd", interestAmtInUSD.String()),
		sdk.NewAttribute("insurance_amount", insuranceAmt.String()),
		sdk.NewAttribute("funding_fee_paid_custody", closedMTP.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", closedMTP.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("borrow_interest_paid_custody", closedMTP.BorrowInterestPaidCustody.String()),
		sdk.NewAttribute("trading_asset_price", tradingAssetPrice.String()),
		sdk.NewAttribute("all_interests_paid", strconv.FormatBool(allInterestsPaid)), // Funding Fee is fully paid but interest amount is only partially paid then this will be false
		sdk.NewAttribute("trigger", trigger),
		sdk.NewAttribute("open_price", closedMTP.OpenPrice.String()),
		sdk.NewAttribute(types.AttributeKeyPerpFee, perpFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeySlippage, slippageFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeKeyWeightBreakingFee, weightBreakingFeesInUsd.String()),
		sdk.NewAttribute(types.AttributeTakerFees, takerFeesInUsd.String()),
	))
}
