package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPConsolidateLiability(mtp *MTP) {
	if mtp.SumCollateral.IsZero() {
		return
	}

	leverage := mtp.Liabilities.Quo(mtp.SumCollateral)
	mtp.ConsolidateLeverage = sdk.NewDecFromInt(leverage)
}
