package keeper

import (
	errorsmod "cosmossdk.io/errors"
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// CalcMinCollateral calculates the minimum collateral required to open a position
func (k Keeper) CalcMinCollateral(ctx sdk.Context, leverage sdkmath.LegacyDec, price sdkmath.LegacyDec, decimals uint64) (sdkmath.Int, error) {
	// leverage must be greater than 1
	if leverage.LTE(sdkmath.LegacyNewDec(1)) {
		return sdkmath.ZeroInt(), errorsmod.Wrapf(types.ErrInvalidLeverage, "leverage must be greater than 1")
	}

	// get min borrow rate
	borrowInterestRateMin := k.GetBorrowInterestRateMin(ctx)

	// round up price
	price = price.Ceil()

	// min_collateral = [ trading_asset_rate_in_usdc / (( leverage - 1 ) * borrow_interest_rate_min ) ] + 10 ^ collateral_decimals
	minCollateral := price.Quo(
		leverage.Sub(sdkmath.LegacyNewDec(1)).Mul(borrowInterestRateMin),
	).Add(sdkmath.LegacyMustNewDecFromStr("10").Power(decimals))

	return minCollateral.TruncateInt(), nil
}
