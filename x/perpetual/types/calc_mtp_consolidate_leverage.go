package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPConsolidateLiability(mtp *MTP) sdk.Dec {
	if mtp.SumCollateral.IsZero() {
		return mtp.ConsolidateLeverage
	}

	leverage := mtp.Liabilities.Quo(mtp.SumCollateral)
	return sdk.NewDecFromInt(leverage)
}
