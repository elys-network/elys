package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/stablestake/types"
)

func (k msgServer) Unbond(goCtx context.Context, msg *types.MsgUnbond) (*types.MsgUnbondResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	params := k.GetParams(ctx)
	creator := sdk.MustAccAddressFromBech32(msg.Creator)
	redemptionRate := k.GetRedemptionRate(ctx)

	shareDenom := types.GetShareDenom()

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

	amountAfterRedemption := params.TotalValue.Sub(redemptionAmount)
	maxAllowed := (params.TotalValue.ToLegacyDec().Mul(params.MaxWithdrawRatio)).TruncateInt()
	if amountAfterRedemption.LT(maxAllowed) {
		return nil, types.ErrInvalidWithdraw
	}

	depositDenom := k.GetDepositDenom(ctx)
	redemptionCoin := sdk.NewCoin(depositDenom, redemptionAmount)
	err = k.bk.SendCoinsFromModuleToAccount(ctx, types.ModuleName, creator, sdk.Coins{redemptionCoin})
	if err != nil {
		return nil, err
	}

	params.TotalValue = params.TotalValue.Sub(redemptionAmount)
	k.SetParams(ctx, params)

	if k.hooks != nil {
		err = k.hooks.AfterUnbond(ctx, creator, msg.Amount)
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
