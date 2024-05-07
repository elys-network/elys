package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Unbond(goCtx context.Context, msg *types.MsgUnbond) (*types.MsgUnbondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	sender := sdk.MustAccAddressFromBech32(msg.Creator)

	shareDenom := types.GetShareDenom()

	// Withdraw committed LP tokens
	err := k.commitmentKeeper.UncommitTokens(ctx, sender, shareDenom, msg.Amount)
	if err != nil {
		return nil, err
	}

	shareCoin := sdk.NewCoin(shareDenom, msg.Amount)
	shareCoins := sdk.NewCoins(shareCoin)
	err = k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.BurnCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	redemptionAmount := sdk.NewDecFromInt(shareCoin.Amount).Mul(params.RedemptionRate).RoundInt()

	depositDenom := k.GetDepositDenom(ctx)
	redemptionCoin := sdk.NewCoin(depositDenom, redemptionAmount)
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{redemptionCoin})
	if err != nil {
		return nil, err
	}

	params.TotalValue = params.TotalValue.Sub(redemptionAmount)
	k.SetParams(ctx, params)

	if k.hooks != nil {
		err := k.hooks.AfterUnbond(ctx, msg.Creator, msg.Amount)
		if err != nil {
			return nil, err
		}
	}

	return &types.MsgUnbondResponse{}, nil
}
