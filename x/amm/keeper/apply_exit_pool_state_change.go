package keeper

import (
	sdkmath "cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
)

func (k Keeper) ApplyExitPoolStateChange(ctx sdk.Context, pool types.Pool, exiter sdk.AccAddress, numShares sdkmath.Int, exitCoins sdk.Coins, isLiquidation bool) error {
	// Withdraw exit amount of token from commitment module to exiter's wallet.
	poolShareDenom := types.GetPoolShareDenom(pool.GetPoolId())

	// Withdraw committed LP tokens
	err := k.commitmentKeeper.UncommitTokens(ctx, exiter, poolShareDenom, numShares, isLiquidation)
	if err != nil {
		return err
	}

	if err = k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()), exiter, exitCoins); err != nil {
		return err
	}

	if err = k.BurnPoolShareFromAccount(ctx, pool, exiter, numShares); err != nil {
		return err
	}

	k.SetPool(ctx, pool)

	types.EmitRemoveLiquidityEvent(ctx, exiter, pool.GetPoolId(), exitCoins)
	if k.hooks != nil {
		err = k.hooks.AfterExitPool(ctx, exiter, pool, numShares, exitCoins)
		if err != nil {
			return err
		}
	}
	return nil
}
