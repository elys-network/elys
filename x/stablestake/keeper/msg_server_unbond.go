package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Unbond(goCtx context.Context, msg *types.MsgUnbond) (*types.MsgUnbondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	pool, found := k.GetPool(ctx, msg.PoolId)
	if !found {
		return nil, types.ErrPoolNotFound
	}

	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	redemptionRate := k.GetRedemptionRateForPool(ctx, pool)

	shareDenom := types.GetShareDenomForPool(pool.PoolId)

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

	amountAfterRedemption := pool.TotalValue.Sub(redemptionAmount)
	maxAllowed := (pool.TotalValue.ToLegacyDec().Mul(pool.MaxLeverageRatio)).TruncateInt()
	if amountAfterRedemption.LT(maxAllowed) {
		return nil, types.ErrInvalidWithdraw
	}

	depositDenom := pool.GetDepositDenom()
	redemptionCoin := sdk.NewCoin(depositDenom, redemptionAmount)
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.Coins{redemptionCoin})
	if err != nil {
		return nil, err
	}

	pool.TotalValue = pool.TotalValue.Sub(redemptionAmount)
	k.SetPool(ctx, pool)

	if k.hooks != nil {
		err = k.hooks.AfterUnbond(ctx, creator, msg.Amount, pool.PoolId)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgUnbondResponse{}, nil
}
