package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"errors"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) PlaceLimitOrder(goCtx context.Context, msg *types.MsgPlaceLimitOrder) (*types.MsgPlaceLimitOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := k.GetPerpetualMarket(ctx, msg.MarketId)
	if err != nil {
		return nil, err
	}

	if err = market.ValidateMsgOpenPosition(*msg); err != nil {
		return nil, err
	}

	_, err = k.GetSubAccount(ctx, sdk.MustAccAddressFromBech32(msg.Creator), market.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "subaccount id: %d", market.Id)
	}

	order := types.PerpetualOrder{
		MarketId:    market.Id,
		OrderType:   msg.OrderType,
		Price:       msg.Price,
		BlockHeight: uint64(ctx.BlockHeight()),
		Owner:       msg.Creator,
		Amount:      msg.BaseQuantity,
	}
	k.SetPerpetualOrder(ctx, order)
	return &types.MsgPlaceLimitOrderResponse{}, nil
}

func (k Keeper) PlaceMarketOrder(goCtx context.Context, msg *types.MsgPlaceMarketOrder) (*types.MsgPlaceMarketOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := k.GetPerpetualMarket(ctx, msg.MarketId)
	if err != nil {
		return nil, err
	}

	//if err = market.ValidateMsgOpenPosition(*msg); err != nil {
	//	return nil, err
	//}

	_, err = k.GetSubAccount(ctx, sdk.MustAccAddressFromBech32(msg.Creator), market.Id)
	if err != nil {
		return nil, errorsmod.Wrapf(err, "subaccount id: %d", market.Id)
	}

	fullyFilled := false
	switch msg.OrderType {
	case types.OrderType_ORDER_TYPE_MARKET_BUY:
		fullyFilled, err = k.ExecuteMarketBuyOrder(ctx, market, *msg)
	case types.OrderType_ORDER_TYPE_MARKET_SELL:
		fullyFilled, err = k.ExecuteMarketSellOrder(ctx, market, *msg)
	default:
		return nil, errorsmod.Wrapf(err, "unknown order type: %s", msg.OrderType)
	}
	if err != nil {
		return nil, err
	}
	if !fullyFilled {
		return nil, errors.New("market order cannot be fully filled")
	}

	return &types.MsgPlaceMarketOrderResponse{}, nil
}
