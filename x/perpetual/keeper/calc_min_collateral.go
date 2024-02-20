package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/perpetual/types"
)

// CalcMinCollateral calculates the minimum collateral required to open a position
func (k Keeper) CalcMinCollateral(ctx sdk.Context, leverage sdk.Dec, price sdk.Dec, decimals uint64) (math.Int, error) {
	// leverage must be greater than 1
	if leverage.LTE(sdk.NewDec(1)) {
		return sdk.ZeroInt(), errorsmod.Wrapf(types.ErrInvalidLeverage, "leverage must be greater than 1")
	}

	// get min borrow rate
	borrowInterestRateMin := k.GetBorrowInterestRateMin(ctx)

	// round up price
	price = price.Ceil()

	// min_collateral = [ trading_asset_rate_in_usdc / (( leverage - 1 ) * borrow_interest_rate_min ) ] + 10 ^ collateral_decimals
	minCollateral := price.Quo(
		leverage.Sub(sdk.NewDec(1)).Mul(borrowInterestRateMin),
	).Add(sdk.MustNewDecFromStr("10").Power(decimals))

	return minCollateral.TruncateInt(), nil
}
