package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	//check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		currentHeight := ctx.BlockHeight()
		_ = currentHeight
		pools := k.GetAllPools(ctx)
		for _, pool := range pools {
			// TODO: fields missing
			// pool.BlockInterestExternal = sdk.ZeroUint()
			// pool.BlockInterestNative = sdk.ZeroUint()
			ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, "")
			if err != nil {
				ctx.Logger().Error(errors.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
				continue // ?
			}
			if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
				rate, err := k.InterestRateComputation(ctx, pool, ammPool)
				if err != nil {
					ctx.Logger().Error(err.Error())
					continue // ?
				}
				pool.InterestRate = rate
				// TODO: field missing
				// pool.LastHeightInterestRateComputed = currentHeight
				_ = k.UpdatePoolHealth(ctx, &pool)
				// TODO: function missing
				// k.TrackSQBeginBlock(ctx, pool)
				mtps, _, _ := k.GetMTPsForPool(ctx, pool.AmmPoolId, nil)
				for _, mtp := range mtps {
					BeginBlockerProcessMTP(ctx, k, mtp, pool, ammPool)
				}
			}
			k.SetPool(ctx, pool)
		}
	}

}
