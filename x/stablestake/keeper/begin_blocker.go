package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

var numBlocks = 15768000 // Number of blocks in 2 year assuming block time 4 seconds

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	k.UpdateInterestForAllPools(ctx)
	// Remove old data, should keep data of 2 years
	if numBlocks < int(ctx.BlockHeight()) {
		delBlock := ctx.BlockHeight() - int64(numBlocks)
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			k.DeleteInterestForPool(ctx, delBlock, pool.Id)
		}
	}
}

func (k Keeper) UpdateInterestForAllPools(ctx sdk.Context) {
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		pool.InterestRate = k.InterestRateComputationForPool(ctx, pool).Dec()
		k.SetPool(ctx, pool)
		k.SetInterestForPool(ctx, types.InterestBlock{
			InterestRate: pool.InterestRate,
			BlockTime:    ctx.BlockTime().Unix(),
			BlockHeight:  uint64(ctx.BlockHeight()),
			PoolId:       pool.Id,
		})
	}
}
