package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ForceCloseShort(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, takeFundPayment bool, baseCurrency string) (math.Int, error) {
	repayAmount := math.ZeroInt()
	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId)
	if err != nil {
		return math.ZeroInt(), err
	}

	// Estimate swap and repay
	repayAmt, err := k.EstimateAndRepay(ctx, mtp, pool, &ammPool, baseCurrency, math.LegacyOneDec())
	if err != nil {
		return math.ZeroInt(), err
	}

	repayAmount = repayAmount.Add(repayAmt)

	address := sdk.MustAccAddressFromBech32(mtp.Address)
	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		params := k.GetParams(ctx)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, *pool, address, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return math.Int{}, err
		}
	}

	return repayAmount, nil
}
