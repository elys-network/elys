package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	currentHeight := ctx.BlockHeight()
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
			rate, err := k.BorrowInterestRateComputation(ctx, pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
				continue
			}
			pool.BorrowInterestRate = rate
			pool.LastHeightBorrowInterestRateComputed = currentHeight

			k.SetBorrowRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.InterestBlock{
				InterestRate: rate,
				BlockTime:    ctx.BlockTime().Unix(),
			})

			err = k.UpdatePoolHealth(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
			err = k.UpdateFundingRate(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}

			k.SetFundingRate(ctx, uint64(ctx.BlockHeight()), pool.AmmPoolId, types.FundingRateBlock{
				FundingRate: pool.FundingRate,
				BlockTime:   ctx.BlockTime().Unix(),
			})
		}
		k.SetPool(ctx, pool)
	}
}
