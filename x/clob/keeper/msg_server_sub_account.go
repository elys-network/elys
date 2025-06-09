package keeper

import (
	"context"
	"fmt"
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	sender := sdk.MustAccAddressFromBech32(msg.Sender)

	crossSubAccount, err := k.GetSubAccount(ctx, sender, types.CrossMarginSubAccountId)
	if err != nil {
		return nil, err
	}

	balance := k.GetSubAccountBalanceOf(ctx, crossSubAccount, msg.Coin.Denom)
	if !balance.IsPositive() {
		return nil, fmt.Errorf("insufficient funds in subaccount for %s", msg.Coin.Denom)
	}

	//check all isolated subaccounts, if they have any residual amount and no open orders, if so send to cross margin account first
	allSubAccounts := k.GetAllOwnerSubAccount(ctx, sender)
	for _, subAccount := range allSubAccounts {
		if subAccount.IsIsolated() {
			minimumBalance, err := k.RequiredMinimumBalance(ctx, subAccount)
			if err != nil {
				return nil, err
			}
			balance := k.GetSubAccountBalance(ctx, subAccount)

			for _, coin := range balance {
				residual := coin.Amount.Sub(minimumBalance.AmountOf(coin.Denom))
				if residual.IsPositive() {
					err = k.TransferFromSubAccountToSubAccount(ctx, subAccount, crossSubAccount, sdk.NewCoins(sdk.NewCoin(coin.Denom, residual)))
					if err != nil {
						return nil, err
					}
				}
			}
		}
	}

	requiredMinimumBalance, err := k.RequiredMinimumBalance(ctx, crossSubAccount)
	if err != nil {
		return nil, err
	}

	withdrawAbleAmount := balance.Amount.Sub(requiredMinimumBalance.AmountOf(msg.Coin.Denom))
	if !withdrawAbleAmount.IsPositive() {
		return nil, fmt.Errorf("insufficient funds in subaccount for %s, max that can be withdrawa: %s", msg.Coin.Denom, withdrawAbleAmount.String())
	}

	if msg.Coin.Amount.GT(withdrawAbleAmount) {
		return nil, fmt.Errorf("withdrawing more than allowed, max withdraw: %s", withdrawAbleAmount.String())
	}
	err = k.SendFromSubAccount(ctx, crossSubAccount, sender, sdk.NewCoins(msg.Coin))
	if err != nil {
		return nil, err
	}
	return &types.MsgWithdrawResponse{}, nil
}
