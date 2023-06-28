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

	// Withdraw token message
	msgWithdrawToken := &ctypes.MsgWithdrawTokens{
		Creator: exiter.String(),
		Amount:  numShares,
		Denom:   poolShareDenom,
	}

	// Withdraw committed LP token
	_, err := msgServer.WithdrawTokens(sdk.WrapSDKContext(ctx), msgWithdrawToken)
	if err != nil {
		return err
	}

	if err := k.bankKeeper.SendCoins(ctx, sdk.MustAccAddressFromBech32(pool.GetAddress()), exiter, exitCoins); err != nil {
		return err
	}

	exitFeeCoins := PortionCoins(exitCoins, pool.PoolParams.ExitFee)
	rebalanceTreasury := sdk.MustAccAddressFromBech32(pool.GetRebalanceTreasury())
	if err := k.bankKeeper.SendCoins(ctx, exiter, rebalanceTreasury, exitFeeCoins); err != nil {
		return err
	}

	if err := k.OnCollectFee(ctx, pool, exitFeeCoins); err != nil {
		return err
	}

	if err := k.BurnPoolShareFromAccount(ctx, pool, exiter, numShares); err != nil {
		return err
	}

	if err := k.SetPool(ctx, pool); err != nil {
		return err
	}

	types.EmitRemoveLiquidityEvent(ctx, exiter, pool.GetPoolId(), exitCoins)
	k.hooks.AfterExitPool(ctx, exiter, pool.GetPoolId(), numShares, exitCoins)
	k.RecordTotalLiquidityDecrease(ctx, exitCoins)
	return nil
}
