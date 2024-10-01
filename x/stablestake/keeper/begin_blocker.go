package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

var numBlocks = 15768000 // Number of blocks in 2 year assuming block time 4 seconds

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)
	params := k.GetParams(ctx)

	if epochPosition == 0 { // if epoch has passed
		rate := k.InterestRateComputation(ctx)
		params.InterestRate = rate

		params.RedemptionRate = k.GetRedemptionRate(ctx)
		k.SetParams(ctx, params)
	}
	k.SetInterest(ctx, uint64(ctx.BlockHeight()), types.InterestBlock{InterestRate: params.InterestRate, BlockTime: ctx.BlockTime().Unix(), BlockHeight: uint64(ctx.BlockHeight())})

	// Remove old data, should keep data of 2 years
	if numBlocks < int(ctx.BlockHeight()) {
		delBlock := ctx.BlockHeight() - int64(numBlocks)
		k.DeleteInterest(ctx, delBlock)
	}
}
