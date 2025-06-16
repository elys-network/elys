package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) ForceClose(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool) (math.Int, math.Int, error) {
	// Estimate swap and repay
	repayAmt, returnAmount, err := k.EstimateAndRepay(ctx, mtp, pool, ammPool, math.LegacyOneDec())
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), err
	}

	address := sdk.MustAccAddressFromBech32(mtp.Address)
	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionClosed(ctx, *ammPool, *pool, address)
		if err != nil {
			return math.Int{}, math.Int{}, err
		}
	}

	return repayAmt, returnAmount, nil
}
