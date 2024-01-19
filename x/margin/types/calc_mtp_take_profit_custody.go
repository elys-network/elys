package types

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPTakeProfitCustody(mtp *MTP) math.Int {
	if IsTakeProfitPriceInifite(mtp) || mtp.TakeProfitPrice.IsZero() {
		return mtp.Custody
	}
	return sdk.NewDecFromInt(mtp.Liabilities).Quo(mtp.TakeProfitPrice).TruncateInt()
}
