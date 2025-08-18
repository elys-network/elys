package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ammtypes "github.com/elys-network/elys/v7/x/amm/types"
	"github.com/elys-network/elys/v7/x/perpetual/types"
)

// ForceClose Three possible cases ForceClose can be called:
// 1. Liquidation - if first we close only 50%
// 2. Stop Loss or Take Profit: We close fully
func (k Keeper) ForceClose(ctx sdk.Context, mtp *types.MTP, pool *types.Pool, ammPool *ammtypes.Pool, isLiquidation bool) (math.Int, math.Int, types.PerpetualFees, math.LegacyDec, math.LegacyDec, error) {
	// Estimate swap and repay
	closingRatio := math.LegacyOneDec()
	addCollateral := false
	if isLiquidation && !mtp.PartialLiquidationDone {
		closingRatio = math.LegacyOneDec().QuoInt64(2)
		addCollateral = true
	}
	repayAmt, returnAmount, perpetualFeesCoins, closingPrice, collateralToAdd, err := k.EstimateAndRepay(ctx, mtp, pool, ammPool, closingRatio, isLiquidation)
	if err != nil {
		return math.ZeroInt(), math.ZeroInt(), types.NewPerpetualFeesWithEmptyCoins(), math.LegacyZeroDec(), math.LegacyZeroDec(), err
	}

	// EpochHooks after perpetual position closed
	if k.hooks != nil {
		address := sdk.MustAccAddressFromBech32(mtp.Address)
		err = k.hooks.AfterPerpetualPositionClosed(ctx, *ammPool, *pool, address, closingRatio, mtp.Id)
		if err != nil {
			return math.Int{}, math.Int{}, types.PerpetualFees{}, math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}

	if addCollateral && !collateralToAdd.IsNil() && collateralToAdd.IsPositive() {
		_, err = k.AddCollateral(ctx, mtp, pool, collateralToAdd, ammPool)
		if err != nil {
			return math.Int{}, math.Int{}, types.PerpetualFees{}, math.LegacyZeroDec(), math.LegacyZeroDec(), err
		}
	}
	return repayAmt, returnAmount, perpetualFeesCoins, closingPrice, closingRatio, nil
}
