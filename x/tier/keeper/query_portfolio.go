package keeper

import (
	"context"

	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/v6/x/tier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PortfolioAll(goCtx context.Context, req *types.QueryAllPortfolioRequest) (*types.QueryAllPortfolioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var portfolios []types.Portfolio
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	portfolioStore := prefix.NewStore(store, types.PortfolioKeyPrefix)

	pageRes, err := query.Paginate(portfolioStore, req.Pagination, func(key []byte, value []byte) error {
		var portfolio types.Portfolio
		if err := k.cdc.Unmarshal(value, &portfolio); err != nil {
			return err
		}

		portfolios = append(portfolios, portfolio)
		return nil
	})

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.QueryAllPortfolioResponse{Portfolio: portfolios, Pagination: pageRes}, nil
}

func (k Keeper) Portfolio(goCtx context.Context, req *types.QueryGetPortfolioRequest) (*types.QueryGetPortfolioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}
	ctx := sdk.UnwrapSDKContext(goCtx)

	user, err := sdk.AccAddressFromBech32(req.User)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	date := k.GetDateFromContext(ctx)
	val, found := k.GetPortfolio(ctx, user, date)
	if !found {
		return nil, status.Error(codes.NotFound, "not found")
	}

	return &types.QueryGetPortfolioResponse{TotalPortfolio: val.Dec().String()}, nil
}
