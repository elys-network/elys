package keeper

import (
	"context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) Deposit(goCtx context.Context, msg *types.MsgDeposit) (*types.MsgDepositResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	sender := sdk.MustAccAddressFromBech32(msg.Sender)
	subAccount, err := k.GetSubAccount(ctx, sender, msg.MarketId)
	if err != nil {
		subAccount = types.SubAccount{
			Owner:            msg.Sender,
			AvailableBalance: sdk.Coins{},
			TotalBalance:     sdk.Coins{},
			TradeNounce:      0,
		}
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
