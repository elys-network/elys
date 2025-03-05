package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) CreateLimitOrder(goCtx context.Context, msg *types.MsgCreateLimitOrder) (*types.MsgCreateLimitOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := k.GetPerpetualMarket(ctx, msg.MarketId)
	if err != nil {
		return nil, err
	}

	if err = market.ValidateMsgOpenPosition(*msg); err != nil {
		return nil, err
	}

	subAccount, err := k.GetSubAccount(ctx, sdk.MustAccAddressFromBech32(msg.Creator), msg.SubAccountId)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "subaccount id: %d", msg.SubAccountId)
	}

	if msg.PerpetualId > 0 {
		perpetual, err := k.GetPerpetual(ctx, msg.MarketId, msg.PerpetualId)
		if err != nil {
			return nil, err
		}
		if perpetual.Owner != msg.Creator || perpetual.SubAccountId != msg.SubAccountId {
			return nil, errors.New("not perpetual owner or subaccount does not holds it")
		}
	}

	collateralAmount := msg.Collateral
	err = k.SendFromSubAccount(ctx, subAccount, market.GetAccount(), sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, collateralAmount)))
	if err != nil {
		return nil, err
	}

	order := types.PerpetualOrder{
		MarketId:     market.Id,
		OrderType:    msg.OrderType,
		Price:        msg.Price,
		BlockHeight:  uint64(ctx.BlockHeight()),
		Owner:        msg.Creator,
		SubAccountId: msg.SubAccountId,
		PerpetualId:  msg.PerpetualId,
		Leverage:     msg.Leverage,
		Collateral:   msg.Collateral,
	}
	k.SetPerpetualOrder(ctx, order)

	return &types.MsgCreateLimitOrderResponse{}, nil
}
