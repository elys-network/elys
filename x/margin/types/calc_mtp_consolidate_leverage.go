package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPConsolidateLiability(mtp *MTP) sdk.Dec {
	if mtp.SumCollateral.IsZero() {
		return sdk.ZeroDec()
	}

	leverage := mtp.Liabilities.Quo(mtp.SumCollateral)
	return sdk.NewDecFromInt(leverage)
}
