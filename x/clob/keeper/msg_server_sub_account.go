package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

// Deposit We always deposit to the cross margin account
// Then for isolated we move balance from margin account to isolated account but before that we check if
// cross margin account has enough balance to cover its own positions and open orders (initial margin + trading fees)
func (k Keeper) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	subAccount, err := k.GetSubAccount(ctx, sender, types.CrossMarginSubAccountId)
	if err != nil {
		subAccount = types.SubAccount{
			Owner:       msg.Sender,
			Id:          types.CrossMarginSubAccountId,
			TradeNounce: 0,
		}
		k.SetSubAccount(ctx, subAccount)
	}

	err = k.AddToSubAccount(ctx, sender, subAccount, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, err
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			types.EventDeposit,
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Coin.String()),
		),
	})

	return &types.MsgDepositResponse{}, nil
}

func (k Keeper) Withdraw(goCtx context.Context, msg *types.MsgWithdraw) (*types.MsgWithdrawResponse, error) {
	//ctx := sdk.UnwrapSDKContext(goCtx)
	//
	//sender := sdk.MustAccAddressFromBech32(msg.Sender)
	//subAccount, err := k.GetSubAccount(ctx, sender, msg.SubAccountId)
	//if err != nil {
	//	return nil, err
	//}

	// 1. Check all open orders
	// 2. Calculate Maximum margin amount + trading fees (maker/taker)
	// 3. Check for crossed margin account and maximum balance to maintain

	return &types.MsgWithdrawResponse{}, nil
}
