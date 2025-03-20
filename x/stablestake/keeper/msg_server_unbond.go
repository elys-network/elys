package keeper

import (
	"context"

	errorsmod "cosmossdk.io/errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Unbond(goCtx context.Context, msg *types.MsgUnbond) (*types.MsgUnbondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	redemptionRate := k.CalculateRedemptionRateForPool(ctx, pool)

	shareDenom := types.GetShareDenomForPool(pool.Id)

	// Withdraw committed LP tokens
	err := k.commitmentKeeper.UncommitTokens(ctx, creator, shareDenom, msg.Amount, false)
	if err != nil {
		return nil, err
	}

	shareCoin := sdk.NewCoin(shareDenom, msg.Amount)
	shareCoins := sdk.NewCoins(shareCoin)
	err = k.bk.SendCoinsFromAccountToModule(ctx, creator, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.BurnCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	redemptionAmount := shareCoin.Amount.ToLegacyDec().Mul(redemptionRate).RoundInt()

	moduleAddr := authtypes.NewModuleAddress(types.ModuleName)
	depositDenom := pool.GetDepositDenom()
	balance := k.bk.GetBalance(ctx, moduleAddr, depositDenom)
	borrowed := pool.TotalValue.Sub(balance.Amount)
	borrowedRatio := (borrowed.ToLegacyDec().Quo(pool.TotalValue.Sub(redemptionAmount).ToLegacyDec()))
	if borrowedRatio.GT(pool.MaxWithdrawRatio) {
		return nil, errorsmod.Wrapf(types.ErrInvalidWithdraw, "borrowedRatio: %d", borrowedRatio)
	}

	redemptionCoin := sdk.NewCoin(depositDenom, redemptionAmount)
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.Coins{redemptionCoin})
	if err != nil {
		return nil, err
	}

	pool.TotalValue = pool.TotalValue.Sub(redemptionAmount)
	k.SetPool(ctx, pool)

	if k.hooks != nil {
		err = k.hooks.AfterUnbond(ctx, creator, msg.Amount, pool.Id)
		if err != nil {
			return nil, err
		}
	}

	ctx.EventManager().EmitEvent(sdk.NewEvent(types.EventUnbond,
		sdk.NewAttribute("address", msg.Creator),
		sdk.NewAttribute("amount", msg.Amount.String()),
		sdk.NewAttribute("shares_burnt", shareCoin.String()),
		sdk.NewAttribute("redemption", redemptionCoin.String()),
	))

	return &types.MsgUnbondResponse{}, nil
}
