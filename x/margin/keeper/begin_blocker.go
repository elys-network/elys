package keeper

import (
	"fmt"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	// check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	// if epoch has not passed
	if epochPosition != 0 {
		return
	}

	// if epoch has passed
	entry, found := k.assetProfileKeeper.GetEntry(ctx, ptypes.BaseCurrency)
	if !found {
		ctx.Logger().Error(errorsmod.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency).Error())
	}
	baseCurrency := entry.Denom

	currentHeight := ctx.BlockHeight()
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, "")
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
			continue // ?
		}
		if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
			rate, err := k.BorrowInterestRateComputation(ctx, pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
				continue // ?
			}
			pool.BorrowInterestRate = rate
			pool.LastHeightBorrowInterestRateComputed = currentHeight
			_ = k.UpdatePoolHealth(ctx, &pool)
			_ = k.UpdateFundingRate(ctx, &pool)
			// TODO: function missing
			// k.TrackSQBeginBlock(ctx, pool)
			mtps, _, _ := k.GetMTPsForPool(ctx, pool.AmmPoolId, nil)
			for _, mtp := range mtps {
				err := BeginBlockerProcessMTP(ctx, k, mtp, pool, ammPool, baseCurrency)
				if err != nil {
					ctx.Logger().Error(err.Error())
					continue // ?
				}
			}
			_ = k.HandleFundingFeeDistribution(ctx, mtps, &pool, ammPool, baseCurrency)
		}
		k.SetPool(ctx, pool)
	}
}
