package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

func (k Keeper) ForceCloseLong(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, takeFundPayment bool, baseCurrency string) (math.Int, error) {
	repayAmount := sdk.ZeroInt()
	// Retrieve AmmPool
	ammPool, err := k.GetAmmPool(ctx, mtp.AmmPoolId, mtp.TradingAsset)
	if err != nil {
		return math.ZeroInt(), err
	}

	err = k.TakeOutCustody(ctx, *mtp, pool, mtp.Custody)
	if err != nil {
		return math.ZeroInt(), err
	}

	// Estimate swap and repay
	repayAmt, err := k.EstimateAndRepay(ctx, *mtp, *pool, ammPool, mtp.Custody, baseCurrency)
	if err != nil {
		return math.ZeroInt(), err
	}

	repayAmount = repayAmount.Add(repayAmt)

	// Hooks after perpetual position closed
	address := sdk.MustAccAddressFromBech32(mtp.Address)
	if k.hooks != nil {
		err := k.hooks.AfterPerpetualPositionClosed(ctx, ammPool, *pool, address)
		if err != nil {
			return math.ZeroInt(), err
		}
	}

	return repayAmount, nil
}
