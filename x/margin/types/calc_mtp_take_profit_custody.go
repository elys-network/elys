package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPTakeProfitCustody(mtp *MTP) sdk.Int {
	return sdk.NewDecFromInt(mtp.Liabilities).Quo(mtp.TakeProfitPrice).TruncateInt()
}
