package keeper

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	assetprofiletypes "github.com/elys-network/elys/x/assetprofile/types"
	ptypes "github.com/elys-network/elys/x/parameter/types"
)

func (k Keeper) BeginBlocker(ctx sdk.Context) {
	//check if epoch has passed then execute
	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition == 0 { // if epoch has passed
		entry, found := k.apKeeper.GetEntry(ctx, ptypes.BaseCurrency)
		if !found {
			ctx.Logger().Error(sdkerrors.Wrapf(assetprofiletypes.ErrAssetProfileNotFound, "asset %s not found", ptypes.BaseCurrency).Error())
		}
		baseCurrency := entry.Denom

		currentHeight := ctx.BlockHeight()
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
				pool.LastHeightInterestRateComputed = currentHeight
				_ = k.UpdatePoolHealth(ctx, &pool)
				// TODO: function missing
				// k.TrackSQBeginBlock(ctx, pool)
				mtps, _, _ := k.GetMTPsForPool(ctx, pool.AmmPoolId, nil)
				for _, mtp := range mtps {
					BeginBlockerProcessMTP(ctx, k, mtp, pool, ammPool, baseCurrency)
				}
			}
			k.SetPool(ctx, pool)
		}
	}

}
