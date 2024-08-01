package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		params := k.GetParams(ctx)
		rate := k.InterestRateComputation(ctx)
		params.InterestRate = rate
		k.SetInterest(ctx, uint64(ctx.BlockHeight()), types.InterestBlock{InterestRate: rate, BlockTime: ctx.BlockTime().Unix()})
		k.SetParams(ctx, params)
	}
}
