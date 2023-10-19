package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Bond(goCtx context.Context, msg *types.MsgBond) (*types.MsgBondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	sender := sdk.MustAccAddressFromBech32(msg.Creator)
	depositCoin := sdk.NewCoin(params.DepositDenom, msg.Amount)
	err := k.bk.SendCoinsFromAccountToModule(ctx, sender, types.ModuleName, sdk.Coins{depositCoin})
	if err != nil {
		return nil, err
	}

	shareDenom := types.GetShareDenom()
	shareAmount := sdk.NewDecFromInt(depositCoin.Amount).Quo(params.RedemptionRate).RoundInt()
	shareCoins := sdk.NewCoins(sdk.NewCoin(shareDenom, shareAmount))

	err = k.bk.MintCoins(ctx, types.ModuleName, shareCoins)
	if err != nil {
		return nil, err
	}

	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, sender, shareCoins)
	if err != nil {
		return nil, err
	}

	return &types.MsgBondResponse{}, nil
}
