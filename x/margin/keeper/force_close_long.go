package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) ForceCloseLong(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, takeFundPayment bool) (sdk.Int, error) {

	// check MTP health against threshold
	safetyFactor := k.GetSafetyFactor(ctx)

	epochLength := k.GetEpochLength(ctx)
	epochPosition := k.GetEpochPosition(ctx, epochLength)

	if epochPosition > 0 {
		repayAmount := sdk.ZeroInt()
		for _, custodyAsset := range mtp.CustodyAssets {
			// Retrieve AmmPool
			ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
			if err != nil {
				return math.ZeroInt(), err
			}

			for _, collateralAsset := range mtp.CollateralAssets {
				// Handle Interest if within epoch position
				if err := k.HandleInterest(ctx, mtp, pool, ammPool, collateralAsset, custodyAsset); err != nil {
					return math.ZeroInt(), err
				}
			}

			if mtp.MtpHealth.GT(safetyFactor) {
				return math.ZeroInt(), types.ErrMTPUnhealthy
			}

			err = k.TakeOutCustody(ctx, *mtp, pool, custodyAsset)
			if err != nil {
				return math.ZeroInt(), err
			}

			for _, collateralAsset := range mtp.CollateralAssets {
				// Estimate swap and repay
				repayAmt, err := k.EstimateAndRepay(ctx, *mtp, *pool, ammPool, collateralAsset, custodyAsset)
				if err != nil {
					return math.ZeroInt(), err
				}

				repayAmount = repayAmount.Add(repayAmt)
			}

			// Hooks after margin position closed
			if k.hooks != nil {
				k.hooks.AfterMarginPositionClosed(ctx, ammPool, *pool)
			}
		}

		return repayAmount, nil
	}

	return math.ZeroInt(), nil
}
