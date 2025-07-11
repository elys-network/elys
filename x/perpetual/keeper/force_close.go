package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v6/x/amm/types"
	"github.com/elys-network/elys/v6/x/perpetual/types"
)

func (k Keeper) ForceClose(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool) (math.Int, math.Int, types.PerpetualFees, math.LegacyDec, error) {
	// Estimate swap and repay
	repayAmt, returnAmount, perpetualFeesCoins, closingPrice, err := k.EstimateAndRepay(ctx, mtp, pool, ammPool, math.LegacyOneDec())
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), types.NewPerpetualFeesWithEmptyCoins(), math.LegacyZeroDec(), err
	}

	address := sdk.MustAccAddressFromBech32(mtp.Address)
	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		err = k.hooks.AfterPerpetualPositionClosed(ctx, *ammPool, *pool, address, math.LegacyOneDec(), mtp.Id)
		if err != nil {
			return math.Int{}, math.Int{}, types.PerpetualFees{}, math.LegacyZeroDec(), err
		}
	}

	return repayAmt, returnAmount, perpetualFeesCoins, closingPrice, nil
}
