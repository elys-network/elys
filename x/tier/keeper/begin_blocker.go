package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	deleteDate := k.GetDateAfterDaysFromContext(ctx, -8)
	// Remove last 100 values at each block
	k.RemovePortfolioLast(ctx, deleteDate, 100)
	// migration does not delete all older entries as we don't know whats last date is there, deleting from past 105 days
	deleteDate = k.GetDateAfterDaysFromContext(ctx, -105)
	k.RemoveLegacyPortfolioCounted(ctx, deleteDate, 100)
}
