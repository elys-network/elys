package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/margin/types"
)

func (k Keeper) ForceCloseLong(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, takeFundPayment bool) (sdk.Int, error) {
	repayAmount := sdk.ZeroInt()
	for _, custody := range mtp.Custodies {
		custodyAsset := custody.Denom
		// Retrieve AmmPool
		ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, custodyAsset)
		if err != nil {
			return math.ZeroInt(), err
		}

		err = k.TakeOutCustody(ctx, *mtp, pool, custodyAsset)
		if err != nil {
			return math.ZeroInt(), err
		}

		for _, collateral := range mtp.Collaterals {
			collateralAsset := collateral.Denom
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
