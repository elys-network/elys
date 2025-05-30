package keeper

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) Close(ctx sdk.Context, msg *types.MsgClose) (*types.MsgCloseResponse, error) {
	closedMtp, repayAmount, closingRatio, returnAmt, fundingFeeAmt, fundingAmtDistributed, interestAmt, insuranceAmt, allInterestsPaid, forceClosed, err := k.ClosePosition(ctx, msg)
	if err != nil {
		return nil, err
	}

	tradingAssetPrice, _, err := k.GetAssetPriceAndAssetUsdcDenomRatio(ctx, closedMtp.TradingAsset)
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(types.EventClose,
			sdk.NewAttribute("mtp_id", strconv.FormatInt(int64(closedMtp.Id), 10)),
			sdk.NewAttribute("owner", closedMtp.Address),
			sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(closedMtp.AmmPoolId), 10)),
			sdk.NewAttribute("collateral_asset", closedMtp.CollateralAsset),
			sdk.NewAttribute("position", closedMtp.Position.String()),
			sdk.NewAttribute("mtp_health", closedMtp.MtpHealth.String()), // should be there if it's partial close
			sdk.NewAttribute("repay_amount", repayAmount.String()),
			sdk.NewAttribute("return_amount", returnAmt.String()),
			sdk.NewAttribute("funding_fee_amount", fundingFeeAmt.String()),
			sdk.NewAttribute("funding_amount_distributed", fundingAmtDistributed.String()),
			sdk.NewAttribute("interest_amount", interestAmt.String()),
			sdk.NewAttribute("insurance_amount", insuranceAmt.String()),
			sdk.NewAttribute("funding_fee_paid_custody", closedMtp.FundingFeePaidCustody.String()),
			sdk.NewAttribute("funding_fee_received_custody", closedMtp.FundingFeeReceivedCustody.String()),
			sdk.NewAttribute("closing_ratio", closingRatio.String()),
			sdk.NewAttribute("trading_asset_price", tradingAssetPrice.String()),
			sdk.NewAttribute("all_interests_paid", strconv.FormatBool(allInterestsPaid)), // Funding Fee is fully paid but interest amount is only partially paid then this will be false
			sdk.NewAttribute("force_closed", strconv.FormatBool(forceClosed)),
		))

	return &types.MsgCloseResponse{
		Id:     closedMtp.Id,
		Amount: repayAmount,
	}, nil
}
