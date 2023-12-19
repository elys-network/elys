package keeper

import (
	"context"
	"errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	commitmentkeeper "github.com/elys-network/elys/x/commitment/keeper"
	ctypes "github.com/elys-network/elys/x/commitment/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Unbond(goCtx context.Context, msg *types.MsgUnbond) (*types.MsgUnbondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	sender := sdk.MustAccAddressFromBech32(msg.Creator)

	shareDenom := types.GetShareDenom()

	// Withdraw committed LP tokens
	msgServer := commitmentkeeper.NewMsgServerImpl(*k.commitmentKeeper)
	_, err := msgServer.UncommitTokens(sdk.WrapSDKContext(ctx), &ctypes.MsgUncommitTokens{
		Creator: sender.String(),
		Amount:  msg.Amount,
		Denom:   shareDenom,
	})
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
	entry, found := k.assetProfileKeeper.GetEntry(ctx, params.DepositDenom)
	if !found {
		return nil, errors.New("invalid denom")
	}

	depositDenom := entry.Denom
	redemptionCoin := sdk.NewCoin(depositDenom, redemptionAmount)
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, sdk.Coins{redemptionCoin})
	if err != nil {
		return nil, err
	}

	params.TotalValue = params.TotalValue.Sub(redemptionAmount)
	k.SetParams(ctx, params)

	return &types.MsgUnbondResponse{}, nil
}
