package keeper

import (
	"context"
	"errors"
	"fmt"

	errorsmod "cosmossdk.io/errors"
	"cosmossdk.io/math"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v6/x/clob/types"
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

	// Additional validations
	if err = k.ValidateOrderPrice(ctx, market, msg.OrderType, msg.Price); err != nil {
		return nil, err
	}

	if err = k.ValidateOrderQuantity(ctx, market, msg.BaseQuantity); err != nil {
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
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	subAccount, err := k.GetSubAccount(ctx, creator, subAccountId)
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
		Filled:       math.LegacyZeroDec(),
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

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventPlaceLimitOrder,
			sdk.NewAttribute(types.AttributeSender, msg.Creator),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", msg.MarketId)),
			sdk.NewAttribute(types.AttributeOrderType, msg.OrderType.String()),
			sdk.NewAttribute(types.AttributePrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeQuantity, msg.BaseQuantity.String()),
			sdk.NewAttribute(types.AttributeOrderId, fmt.Sprintf("%d", order.Counter)),
		),
	)

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

	// Additional validations
	if err = k.ValidateOrderQuantity(ctx, market, msg.BaseQuantity); err != nil {
		return nil, err
	}

	// Check if there's sufficient liquidity for market order
	if msg.OrderType == types.OrderType_ORDER_TYPE_MARKET_BUY {
		sellOrders := k.GetSellOrdersUpToQuantity(ctx, market.Id, msg.BaseQuantity)
		if len(sellOrders) == 0 {
			return nil, types.ErrNoOrdersAvailable.Wrapf("no sell orders available for market buy in market %d", market.Id)
		}
		totalAvailable := math.LegacyZeroDec()
		for _, order := range sellOrders {
			totalAvailable = totalAvailable.Add(order.Quantity)
		}
		if totalAvailable.LT(msg.BaseQuantity) {
			return nil, WrapInsufficientLiquidityError(totalAvailable, msg.BaseQuantity, "market buy")
		}
	} else if msg.OrderType == types.OrderType_ORDER_TYPE_MARKET_SELL {
		buyOrders := k.GetBuyOrdersUpToQuantity(ctx, market.Id, msg.BaseQuantity)
		if len(buyOrders) == 0 {
			return nil, types.ErrNoOrdersAvailable.Wrapf("no buy orders available for market sell in market %d", market.Id)
		}
		totalAvailable := math.LegacyZeroDec()
		for _, order := range buyOrders {
			totalAvailable = totalAvailable.Add(order.Quantity)
		}
		if totalAvailable.LT(msg.BaseQuantity) {
			return nil, WrapInsufficientLiquidityError(totalAvailable, msg.BaseQuantity, "market sell")
		}
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
	creator, err := sdk.AccAddressFromBech32(msg.Creator)
	if err != nil {
		return nil, err
	}
	subAccount, err := k.GetSubAccount(ctx, creator, subAccountId)
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
		return nil, types.ErrOrderNotFilled.Wrap("market order cannot be fully filled")
	}

	// Emit event
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			types.EventPlaceMarketOrder,
			sdk.NewAttribute(types.AttributeSender, msg.Creator),
			sdk.NewAttribute(types.AttributeMarketId, fmt.Sprintf("%d", msg.MarketId)),
			sdk.NewAttribute(types.AttributeOrderType, msg.OrderType.String()),
			sdk.NewAttribute(types.AttributeQuantity, msg.BaseQuantity.String()),
		),
	)

	return &types.MsgPlaceMarketOrderResponse{}, nil
}
