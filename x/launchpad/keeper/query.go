package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/x/launchpad/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Bonus(goCtx context.Context, req *types.QueryBonusRequest) (*types.QueryBonusResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	allOrders := k.GetAllOrders(ctx)
	orders := []types.Purchase{}
	for _, order := range allOrders {
		if order.OrderMaker == req.User {
			orders = append(orders, order)
		}
	}
	bonusAmount := sdk.ZeroInt()
	for _, order := range orders {
		bonusAmount = bonusAmount.Add(order.BonusAmount)
	}

	return &types.QueryBonusResponse{
		TotalBonus: bonusAmount,
	}, nil
}

func (k Keeper) BuyElysEst(goCtx context.Context, req *types.QueryBuyElysEstRequest) (*types.QueryBuyElysEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	elysAmount, orders, err := k.CalcBuyElysResult(ctx, "", req.SpendingToken, req.TokenAmount)
	if err != nil {
		return nil, err
	}
	bonusAmount := sdk.ZeroInt()
	for _, order := range orders {
		bonusAmount = bonusAmount.Add(order.BonusAmount)
	}

	return &types.QueryBuyElysEstResponse{
		ElysAmount:  elysAmount,
		BonusAmount: bonusAmount,
		Orders:      orders,
	}, nil
}

func (k Keeper) ReturnElysEst(goCtx context.Context, req *types.QueryReturnElysEstRequest) (*types.QueryReturnElysEstResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)
	returnTokenAmount, err := k.CalcReturnElysResult(ctx, req.OrderId, req.ElysAmount)
	if err != nil {
		return nil, err
	}

	return &types.QueryReturnElysEstResponse{
		Amount: returnTokenAmount,
	}, nil
}

func (k Keeper) Orders(goCtx context.Context, req *types.QueryOrdersRequest) (*types.QueryOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	allOrders := k.GetAllOrders(ctx)
	orders := []types.Purchase{}
	for _, order := range allOrders {
		if order.OrderMaker == req.User {
			orders = append(orders, order)
		}
	}

	return &types.QueryOrdersResponse{
		Purchases: orders,
	}, nil
}

func (k Keeper) AllOrders(goCtx context.Context, req *types.QueryAllOrdersRequest) (*types.QueryAllOrdersResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	allOrders := k.GetAllOrders(ctx)
	return &types.QueryAllOrdersResponse{
		Purchases: allOrders,
	}, nil
}
