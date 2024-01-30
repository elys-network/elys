package types

import (
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func GenerateOpenEvent(mtp *MTP) sdk.Event {
	return sdk.NewEvent(EventOpen,
		sdk.NewAttribute("id", strconv.FormatInt(int64(mtp.Id), 10)),
		sdk.NewAttribute("position", mtp.Position.String()),
		sdk.NewAttribute("address", mtp.Address),
		sdk.NewAttribute("collateral", mtp.Collateral.String()),
		sdk.NewAttribute("custody", mtp.Custody.String()),
		sdk.NewAttribute("leverage", mtp.Leverage.String()),
		sdk.NewAttribute("liabilities", mtp.Liabilities.String()),
		sdk.NewAttribute("borrow_interest_paid_collateral", mtp.BorrowInterestPaidCollateral.String()),
		sdk.NewAttribute("borrow_interest_paid_custody", mtp.BorrowInterestPaidCustody.String()),
		sdk.NewAttribute("borrow_interest_unpaid_collateral", mtp.BorrowInterestUnpaidCollateral.String()),
		sdk.NewAttribute("health", mtp.MtpHealth.String()),
	)
}
