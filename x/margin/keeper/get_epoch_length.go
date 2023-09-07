package keeper

import sdk "github.com/cosmos/cosmos-sdk/types"

func (k Keeper) GetEpochLength(ctx sdk.Context) int64 {
	return k.GetParams(ctx).EpochLength
}
