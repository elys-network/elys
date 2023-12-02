package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// ReachedTakeProfitPrice tells if the take profit price is reached
func ReachedTakeProfitPrice(mtp *MTP, assetPrice sdk.Int) bool {
	if mtp.Position == Position_LONG {
		return mtp.TakeProfitPrice.GTE(sdk.NewDecFromInt(assetPrice))
	} else if mtp.Position == Position_SHORT {
		return mtp.TakeProfitPrice.LTE(sdk.NewDecFromInt(assetPrice))
	}
	return false
}
