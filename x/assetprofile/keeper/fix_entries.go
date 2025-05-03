package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"strings"
)

func (k Keeper) FixEntries(ctx sdk.Context) {
	allEntries := k.GetAllEntry(ctx)
	var toDelete []string

	for _, entry := range allEntries {
		if strings.HasPrefix(entry.BaseDenom, "ibc/") {
			toDelete = append(toDelete, entry.BaseDenom)
		}
	}
	for _, denom := range toDelete {
		k.RemoveEntry(ctx, denom)
	}
}
