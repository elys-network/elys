package types

import (
	"fmt"
	"strconv"

	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GenerateCloseEvent(mtp *MTP, repayAmount math.Int) sdk.Event {
	return sdk.NewEvent(EventClose,
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral_asset", mtp.CollateralAsset),
		sdk.NewAttribute("trading_asset", mtp.TradingAsset),
		sdk.NewAttribute("liabilities_asset", mtp.LiabilitiesAsset),
		sdk.NewAttribute("custody_asset", mtp.CustodyAsset),
		sdk.NewAttribute("collateral", mtp.Collateral.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("borrow_interest_paid_collateral", mtp.BorrowInterestPaidCollateral.String()),
		sdk.NewAttribute("borrow_interest_paid_custody", mtp.BorrowInterestPaidCustody.String()),
		sdk.NewAttribute("borrow_interest_unpaid_collateral", mtp.BorrowInterestUnpaidCollateral.String()),
		sdk.NewAttribute("custody", mtp.Custody.String()),
		sdk.NewAttribute("take_profit_liabilities", mtp.TakeProfitLiabilities.String()),
		sdk.NewAttribute("take_profit_custody", mtp.TakeProfitCustody.String()),
		sdk.NewAttribute("leverage", fmt.Sprintf("%s", mtp.Leverage)),
		sdk.NewAttribute("mtp_health", mtp.MtpHealth.String()),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("amm_pool_id", strconv.FormatInt(int64(mtp.AmmPoolId), 10)),
		sdk.NewAttribute("consolidate_leverage", fmt.Sprintf("%s", mtp.ConsolidateLeverage)),
		sdk.NewAttribute("sum_collateral", mtp.SumCollateral.String()),
		sdk.NewAttribute("take_profit_price", mtp.TakeProfitPrice.String()),
		sdk.NewAttribute("take_profit_borrow_rate", mtp.TakeProfitBorrowRate.String()),
		sdk.NewAttribute("funding_fee_paid_collateral", mtp.FundingFeePaidCollateral.String()),
		sdk.NewAttribute("funding_fee_paid_custody", mtp.FundingFeePaidCustody.String()),
		sdk.NewAttribute("funding_fee_received_collateral", mtp.FundingFeeReceivedCollateral.String()),
		sdk.NewAttribute("funding_fee_received_custody", mtp.FundingFeeReceivedCustody.String()),
		sdk.NewAttribute("open_price", mtp.OpenPrice.String()),
		sdk.NewAttribute("repay_amount", repayAmount.String()),
	)
}
