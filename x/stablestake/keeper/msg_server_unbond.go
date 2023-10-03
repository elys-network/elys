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
	shareCoin := sdk.NewCoin(shareDenom, msg.Amount)
	shareCoins := sdk.NewCoins(shareCoin)
	err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.BurnCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	redemptionAmount := sdk.NewDecFromInt(shareCoin.Amount).Mul(params.RedemptionRate).RoundInt()
	redemptionCoin := sdk.NewCoin(params.DepositDenom, redemptionAmount)
	err = k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{redemptionCoin})
	if err != nil {
		return nil, err
	}

	return &types.MsgUnbondResponse{}, nil
}
