package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)
	params := k.GetParams(ctx)

	if epochPosition == 0 { // if epoch has passed
		rate := k.InterestRateComputation(ctx)
		params.InterestRate = rate
		k.SetParams(ctx, params)
	}
	// TODO: Remove old data, should keep data of 3 years or we could do last block which is required by debt
	k.SetInterest(ctx, uint64(ctx.BlockHeight()), types.InterestBlock{InterestRate: params.InterestRate, BlockTime: ctx.BlockTime().Unix()})
}
