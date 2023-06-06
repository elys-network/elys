package keeper

import (
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

// swapExactAmountOut is an internal method for swapping an exact amount of tokens
// as input to a pool, using the provided swapFee. This is intended to allow
// different swap fees as determined by multi-hops, or when recovering from
// chain liveness failures.
// TODO: investigate if swapFee can be unexported
// https://github.com/osmosis-labs/osmosis/issues/3130
func (k Keeper) SwapExactAmountOut(
	ctx sdk.Context,
	sender sdk.AccAddress,
	pool types.Pool,
	tokenInDenom string,
	tokenOutMaxAmount math.Int,
	tokenCoin sdk.Coin,
	swapFee sdk.Dec,
) (tokenOutAmount math.Int, err error) {
	return math.Int{}, nil
}
