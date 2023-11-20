package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func CalcMTPTakeProfitCustodies(mtp *MTP) sdk.Coins {
	takeProfitCustodies := mtp.TakeProfitCustodies
	if IsTakeProfitPriceInifite(mtp) {
		return takeProfitCustodies
	}
	for custodyIndex := range mtp.Custodies {
		takeProfitCustodies[custodyIndex].Amount = sdk.NewDecFromInt(mtp.Liabilities).Quo(mtp.TakeProfitPrice).TruncateInt()
	}
	return takeProfitCustodies
}
