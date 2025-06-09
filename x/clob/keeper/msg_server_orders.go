package keeper

import (
	"context"
	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/clob/types"
)

func (k Keeper) PlaceLimitOrder(goCtx context.Context, msg *types.MsgPlaceLimitOrder) (*types.MsgPlaceLimitOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := k.GetPerpetualMarket(ctx, msg.MarketId)
	if err != nil {
		return nil, err
	}

	if err = market.ValidateOpenPositionRequest(msg.MarketId, msg.Price, msg.BaseQuantity, false); err != nil {
		return nil, err
	}

	subAccountId := types.CrossMarginSubAccountId
	if msg.IsIsolated {
		subAccountId = market.Id
	}
	crossMarginAccount := types.SubAccount{
		Owner:       msg.Creator,
		Id:          types.CrossMarginSubAccountId,
		TradeNounce: 0,
	}
	subAccount, err := k.GetSubAccount(ctx, sdk.MustAccAddressFromBech32(msg.Creator), subAccountId)
	if err != nil {
		if errors.Is(err, types.ErrSubAccountNotFound) {
			subAccount = types.SubAccount{
				Owner:       msg.Creator,
				Id:          subAccountId,
				TradeNounce: 0,
			}
		} else {
			return nil, err
		}
	}

	crossMarginRequiredMinimumBalance, err := k.RequiredMinimumBalance(ctx, crossMarginAccount)
	if err != nil {
		return nil, err
	}

	counter := k.GetAndIncrementOrderCounter(ctx, market.Id)
	order := types.PerpetualOrder{
		MarketId:     market.Id,
		OrderType:    msg.OrderType,
		Price:        msg.Price,
		Counter:      counter,
		Owner:        msg.Creator,
		SubAccountId: subAccount.Id,
		Amount:       msg.BaseQuantity,
	}

	orderRequiredBalance, err := k.RequiredBalanceForOrder(ctx, order)
	if err != nil {
		return nil, err
	}

	// freeCollateral = currentBalance - requiredBalance
	freeCollateral := k.GetSubAccountBalanceOf(ctx, crossMarginAccount, market.QuoteDenom).Amount.Sub(crossMarginRequiredMinimumBalance.AmountOf(market.QuoteDenom))

	// Check the minimum balance in the subaccount required to execute the position
	// crossMarginBalance >= crossMarginMinimumBalance + orderRequiredBalance
	if !freeCollateral.GTE(orderRequiredBalance.Amount) {
		return nil, fmt.Errorf("insufficient balance, deposit more. free: %s, required minimum: %s, required by order: %s", freeCollateral.String(), crossMarginRequiredMinimumBalance.AmountOf(market.QuoteDenom).String(), orderRequiredBalance.String())
	}

	if subAccountId != types.CrossMarginSubAccountId {
		err = k.TransferFromSubAccountToSubAccount(ctx, crossMarginAccount, subAccount, sdk.NewCoins(orderRequiredBalance))
		if err != nil {
			return nil, err
		}
	}

	subAccount.TradeNounce++
	k.SetSubAccount(ctx, subAccount)

	k.SetPerpetualOrder(ctx, order)
	k.SetOrderOwner(ctx, types.PerpetualOrderOwner{
		Owner:        msg.Creator,
		SubAccountId: subAccount.Id,
		OrderKey:     types.NewOrderKey(order.MarketId, order.OrderType, order.Price, order.Counter),
	})
	return &types.MsgPlaceLimitOrderResponse{}, nil
}

func (k Keeper) PlaceMarketOrder(goCtx context.Context, msg *types.MsgPlaceMarketOrder) (*types.MsgPlaceMarketOrderResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)

	market, err := k.GetPerpetualMarket(ctx, msg.MarketId)
	if err != nil {
		return nil, err
	}

	if err = market.ValidateOpenPositionRequest(msg.MarketId, math.LegacyZeroDec(), msg.BaseQuantity, true); err != nil {
		return nil, err
	}

	subAccountId := types.CrossMarginSubAccountId
	if msg.IsIsolated {
		subAccountId = market.Id
	}
	crossMarginAccount := types.SubAccount{
		Owner:       msg.Creator,
		Id:          types.CrossMarginSubAccountId,
		TradeNounce: 0,
	}
	subAccount, err := k.GetSubAccount(ctx, sdk.MustAccAddressFromBech32(msg.Creator), subAccountId)
	if err != nil {
		if errors.Is(err, types.ErrSubAccountNotFound) {
			subAccount = types.SubAccount{
				Owner:       msg.Creator,
				Id:          subAccountId,
				TradeNounce: 0,
			}
		} else {
			return nil, err
		}
	}

	crossMarginRequiredMinimumBalance, err := k.RequiredMinimumBalance(ctx, crossMarginAccount)
	if err != nil {
		return nil, err
	}

	crossMarginBalance := k.GetSubAccountBalanceOf(ctx, crossMarginAccount, market.QuoteDenom)

	freeCollateral := crossMarginBalance.Amount.Sub(crossMarginRequiredMinimumBalance.AmountOf(market.QuoteDenom))
	if !freeCollateral.IsPositive() {
		return nil, fmt.Errorf("insufficient balance, balance: %s, required to maintain: %s", crossMarginBalance.String(), crossMarginRequiredMinimumBalance.String())
	}

	if subAccountId != types.CrossMarginSubAccountId {
		err = k.TransferFromSubAccountToSubAccount(ctx, crossMarginAccount, subAccount, sdk.NewCoins(sdk.NewCoin(market.QuoteDenom, freeCollateral)))
		if err != nil {
			return nil, err
		}
	}

	subAccount.TradeNounce++
	k.SetSubAccount(ctx, subAccount)

	fullyFilled := false
	switch msg.OrderType {
	case types.OrderType_ORDER_TYPE_MARKET_BUY:
		fullyFilled, err = k.ExecuteMarketBuyOrder(ctx, market, *msg, false, true)
	case types.OrderType_ORDER_TYPE_MARKET_SELL:
		fullyFilled, err = k.ExecuteMarketSellOrder(ctx, market, *msg, false, false)
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
