package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/amm/types"
)

// ApplyDiscount applies discount to swap fee if applicable
func (k Keeper) ApplyDiscount(ctx sdk.Context, swapFee sdk.Dec, discount sdk.Dec, sender string) (sdk.Dec, sdk.Dec, error) {
	// if discount is nil, return swap fee and zero discount
	if discount.IsNil() {
		return swapFee, sdk.ZeroDec(), nil
	}

	// if discount is zero, return swap fee and zero discount
	if discount.IsZero() {
		return swapFee, sdk.ZeroDec(), nil
	}

	// check if discount is positive and signer address is broker address otherwise throw an error
	if discount.IsPositive() && sender != k.GetBrokerAddress(ctx).String() {
		return sdk.ZeroDec(), sdk.ZeroDec(), sdkerrors.Wrapf(types.ErrInvalidDiscount, "discount %s is positive and signer address %s is not broker address %s", discount, sender, k.GetBrokerAddress(ctx))
	}

	// apply discount percentage to swap fee
	swapFee = swapFee.Mul(sdk.OneDec().Sub(discount))

	return swapFee, discount, nil
}
