package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/amm/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
)

func (k Keeper) applyExitPoolStateChange(ctx sdk.Context, pool types.Pool, exiter sdk.AccAddress, numShares sdk.Int, exitCoins sdk.Coins) error {
	// Withdraw exit amount of token from commitment module to exiter's wallet.
	msgServer := commitmentkeeper.NewMsgServerImpl(*k.commitmentKeeper)

	poolShareDenom := types.GetPoolShareDenom(pool.GetPoolId())
	amount := exitCoins.AmountOf(poolShareDenom)

	// Withdraw token msag
	msgWithdrawToken := &ctypes.MsgWithdrawTokens{
		Creator: exiter.String(),
		Amount:  amount,
		Denom:   poolShareDenom,
	}

	// Withdraw committed LP token
	_, err := msgServer.WithdrawTokens(ctx, msgWithdrawToken)
	if err != nil {
		return err
	}

	err = k.bankKeeper.SendCoins(ctx, sdk.AccAddress(pool.GetAddress()), exiter, exitCoins)
	if err != nil {
		return err
	}

	err = k.BurnPoolShareFromAccount(ctx, pool, exiter, numShares)
	if err != nil {
		return err
	}

	err = k.SetPool(ctx, pool)
	if err != nil {
		return err
	}

	types.EmitRemoveLiquidityEvent(ctx, exiter, pool.GetPoolId(), exitCoins)
	k.hooks.AfterExitPool(ctx, exiter, pool.GetPoolId(), numShares, exitCoins)
	k.RecordTotalLiquidityDecrease(ctx, exitCoins)
	return nil
}
