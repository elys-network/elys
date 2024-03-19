package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// TODO: implement vesting based on time gone
}
