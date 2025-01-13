package keeper

import (
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) EmitForceClose(ctx sdk.Context, trigger string, mtp types.MTP, repayAmount, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt math.Int, closer string, allInterestsPaid bool, tradingAssetPrice math.LegacyDec) {
	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventForceClosed,
		sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("owner", mtp.Address),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("closer", closer),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
		sdk.NewAttribute("return_amount", returnAmt.String()),
		sdk.NewAttribute("funding_fee_amount", fundingFeeAmt.String()),
		sdk.NewAttribute("funding_amount_distributed", fundingAmtDistributed.String()),
		sdk.NewAttribute("interest_amount", interestAmt.String()),
		sdk.NewAttribute("insurance_amount", insuranceAmt.String()),
		sdk.NewAttribute("funding_fee_paid_custody", mtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_custody", mtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("trading_asset_price", tradingAssetPrice.String()),
		sdk.NewAttribute("all_interests_paid", strconv.FormatBool(allInterestsPaid)), // Funding Fee is fully paid but interest amount is only partially paid then this will be false
		sdk.NewAttribute("trigger", trigger),
	))
}
