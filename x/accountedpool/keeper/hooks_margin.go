package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/x/amm/types"
	margintypes "github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) AfterMarginPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPoolByMargin(ctx, ammPool, marginPool)
}

func (k Keeper) AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPoolByMargin(ctx, ammPool, marginPool)
}

func (k Keeper) AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	k.UpdateAccountedPoolByMargin(ctx, ammPool, marginPool)
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

func (h MarginHooks) AfterMarginPositionOpended(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterMarginPositionOpended(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterMarginPositionModified(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterMarginPositionModified(ctx, ammPool, marginPool)
}

func (h MarginHooks) AfterMarginPositionClosed(ctx sdk.Context, ammPool ammtypes.Pool, marginPool margintypes.Pool) {
	h.k.AfterMarginPositionClosed(ctx, ammPool, marginPool)
}
