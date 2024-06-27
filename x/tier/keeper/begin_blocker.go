package keeper

import (
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	year, month, day := ctx.BlockTime().Date()
	dateToday := time.Date(year, month, day, 0, 0, 0, 0, ctx.BlockTime().Location())
	deleteDate := dateToday.AddDate(0, 0, -8)
	// Remove last 500 values at each block
	k.RemovePortfolioLast(ctx, deleteDate.String(), 500)
}
