package types

import sdkmath "cosmossdk.io/math"

func CalcMTPConsolidateLiability(mtp *MTP) sdkmath.LegacyDec {
	if mtp.SumCollateral.IsZero() {
		return mtp.ConsolidateLeverage
	}

	leverage := mtp.Liabilities.Quo(mtp.SumCollateral)
	return sdkmath.LegacyNewDecFromInt(leverage)
}
