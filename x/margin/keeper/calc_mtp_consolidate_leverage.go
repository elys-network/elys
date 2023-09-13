package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) CalcMTPConsolidateLiability(ctx sdk.Context, mtp *types.MTP) {
	if mtp.SumCollateral.IsZero() {
		return
	}

	leverage := mtp.Liabilities.Quo(mtp.SumCollateral)
	mtp.ConsolidateLeverage = sdk.NewDecFromInt(leverage)
}
