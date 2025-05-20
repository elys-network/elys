package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v4/x/amm/types"
	"github.com/elys-network/elys/v4/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

func (k Keeper) ForceClose(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool) (math.Int, math.Int, error) {
	repayAmount := math.ZeroInt()

	// Estimate swap and repay
	repayAmt, returnAmount, err := k.EstimateAndRepay(ctx, mtp, pool, ammPool, osmomath.OneBigDec())
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), err
	}

	repayAmount = repayAmount.Add(repayAmt)

	address := sdk.MustAccAddressFromBech32(mtp.Address)
	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		params := k.GetParams(ctx)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, *ammPool, *pool, address, params.EnableTakeProfitCustodyLiabilities)
		if err != nil {
			return math.Int{}, math.Int{}, err
		}
	}

	return repayAmt, returnAmount, nil
}
