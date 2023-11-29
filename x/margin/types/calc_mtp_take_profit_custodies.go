package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPTakeProfitCustodies(mtp *MTP) sdk.Coins {
	takeProfitCustodies := mtp.TakeProfitCustodies
	for custodyIndex, custody := range mtp.Custodies {
		if IsTakeProfitPriceInifite(mtp) || mtp.TakeProfitPrice.IsZero() {
			takeProfitCustodies[custodyIndex].Amount = custody.Amount
			continue
		}
		takeProfitCustodies[custodyIndex].Amount = sdk.NewDecFromInt(mtp.Liabilities).Quo(mtp.TakeProfitPrice).TruncateInt()
	}
	return takeProfitCustodies
}
