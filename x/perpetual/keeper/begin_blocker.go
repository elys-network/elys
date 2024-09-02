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
	baseCurrencyDecimal := entry.Decimals

	currentHeight := ctx.BlockHeight()
	pools := k.GetAllPools(ctx)
	for _, pool := range pools {
		ammPool, err := k.GetAmmPool(ctx, pool.AmmPoolId, "")
		if err != nil {
			ctx.Logger().Error(errorsmod.Wrap(err, fmt.Sprintf("error getting amm pool: %d", pool.AmmPoolId)).Error())
			continue
		}
		if k.IsPoolEnabled(ctx, pool.AmmPoolId) {
			// TODO: update borrow interest rate, cumulative borrow interest rate store
			rate, err := k.BorrowInterestRateComputation(ctx, pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
				continue
			}
			pool.BorrowInterestRate = rate
			pool.LastHeightBorrowInterestRateComputed = currentHeight
			err = k.UpdatePoolHealth(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
			// TODO: update funding rate, cumulative funding rate store
			err = k.UpdateFundingRate(ctx, &pool)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}

			mtps, _, _ := k.GetMTPsForPool(ctx, pool.AmmPoolId, nil)
			for _, mtp := range mtps {
				err := BeginBlockerProcessMTP(ctx, k, mtp, pool, ammPool, baseCurrency, baseCurrencyDecimal)
				if err != nil {
					ctx.Logger().Error(err.Error())
					continue
				}
			}
			err = k.HandleFundingFeeDistribution(ctx, mtps, &pool, ammPool, baseCurrency)
			if err != nil {
				ctx.Logger().Error(err.Error())
			}
		}
		k.SetPool(ctx, pool)
	}
}
