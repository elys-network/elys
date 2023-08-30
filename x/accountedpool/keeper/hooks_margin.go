package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) AfterMarginPositionOpended(ctx sdk.Context, poolId uint64) {
	k.UpdateAccountedPool(ctx, poolId)
}

func (k Keeper) AfterMarginPositionModified(ctx sdk.Context, poolId uint64) {
	k.UpdateAccountedPool(ctx, poolId)
}

func (k Keeper) AfterMarginPositionClosed(ctx sdk.Context, poolId uint64) {
	k.UpdateAccountedPool(ctx, poolId)
}

// Hooks wrapper struct for tvl keeper
type MarginHooks struct {
	k Keeper
}

var _ margintypes.MarginHooks = MarginHooks{}

// Return the wrapper struct
func (k Keeper) MarginHooks() MarginHooks {
	return MarginHooks{k}
}

func (h MarginHooks) AfterMarginPositionOpended(ctx sdk.Context, poolId uint64) {
	h.k.AfterMarginPositionOpended(ctx, poolId)
}

func (h MarginHooks) AfterMarginPositionModified(ctx sdk.Context, poolId uint64) {
	h.k.AfterMarginPositionModified(ctx, poolId)
}

func (h MarginHooks) AfterMarginPositionClosed(ctx sdk.Context, poolId uint64) {
	h.k.AfterMarginPositionClosed(ctx, poolId)
}
