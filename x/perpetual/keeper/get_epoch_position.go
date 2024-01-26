package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

// get position of current block in epoch
func (k Keeper) GetEpochPosition(ctx sdk.Context, epochLength int64) int64 {
	if epochLength <= 0 {
		epochLength = 1
	}
	currentHeight := ctx.BlockHeight()
	return currentHeight % epochLength
}
