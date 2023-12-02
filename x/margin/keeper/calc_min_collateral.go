package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/margin/types"
)

// CalcMinCollateral calculates the minimum collateral required to open a position
func (k Keeper) CalcMinCollateral(ctx sdk.Context, leverage sdk.Dec) (sdk.Int, error) {
	// leverage must be greater than 1
	if leverage.LTE(sdk.NewDec(1)) {
		return sdk.ZeroInt(), sdkerrors.Wrapf(types.ErrInvalidLeverage, "leverage must be greater than 1")
	}

	// get min borrow rate
	borrowInterestRateMin := k.GetBorrowInterestRateMin(ctx)

	// min_collateral = ( 1 / (( leverage - 1 ) * borrow_interest_rate_min )) / 10 ^ collateral_decimals
	minCollateral := sdk.NewDec(1).Quo(
		leverage.Sub(sdk.NewDec(1)).Mul(borrowInterestRateMin),
	)

	return minCollateral.TruncateInt(), nil
}
