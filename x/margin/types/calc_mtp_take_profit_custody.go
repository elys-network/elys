package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPTakeProfitCustody(mtp *MTP) sdk.Int {
	if IsTakeProfitPriceInifite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return mtp.Custody
	}
	return sdk.NewDecFromInt(mtp.Liabilities).Quo(mtp.TakeProfitPrice).TruncateInt()
}
