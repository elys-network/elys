package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k Keeper) UnwindCoin(ctx sdk.Context, coin sdk.Coin) error {
	// Case 1: coin is a LP share
	// if strings.HasPrefix(coin.Denom, "lp-") {
	// 	// exit pool
	// 	_, _, _, _, _, err := k.amm.ExitPool(ctx, vaultAddress, perform_action.ExitPool.PoolId, perform_action.ExitPool.ShareAmountIn, perform_action.ExitPool.MinAmountsOut, perform_action.ExitPool.TokenOutDenom, false, true)
	// 	if err != nil {
	// 		return errorsmod.Wrapf(types.ErrInvalidAction, "action failed with error: %s", err)
	// 	}
	// }
	return nil
}
