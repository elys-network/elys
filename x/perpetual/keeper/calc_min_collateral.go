package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
	"github.com/osmosis-labs/osmosis/osmomath"
)

// CalcMinCollateral calculates the minimum collateral required to open a position
func (k Keeper) CalcMinCollateral(ctx sdk.Context, leverage osmomath.BigDec, price osmomath.BigDec, decimals uint64) (sdkmath.Int, error) {
	// leverage must be greater than 1
	if leverage.LTE(osmomath.OneBigDec()) {
		return sdkmath.ZeroInt(), errorsmod.Wrapf(types.ErrInvalidLeverage, "leverage must be greater than 1")
	}

	// get min borrow rate
	borrowInterestRateMin := k.GetBigDecBorrowInterestRateMin(ctx)

	// round up price
	price = price.Ceil()

	// min_collateral = [ trading_asset_rate_in_usdc / (( leverage - 1 ) * borrow_interest_rate_min ) ] + 10 ^ collateral_decimals
	minCollateral := price.Quo(
		leverage.Sub(osmomath.OneBigDec()).Mul(borrowInterestRateMin),
	).Add(osmomath.MustNewBigDecFromStr("10").PowerInteger(decimals))

	return minCollateral.Dec().TruncateInt(), nil
}
