package keeper

import (
	"context"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/membershiptier/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) PortfolioAll(goCtx context.Context, req *types.QueryAllPortfolioRequest) (*types.QueryAllPortfolioResponse, error) {
	if req == nil {
		return nil, status.Error(codes.InvalidArgument, "invalid request")
	}

	var portfolios []types.Portfolio
	ctx := sdk.UnwrapSDKContext(goCtx)

	store := ctx.KVStore(k.storeKey)
	portfolioStore := prefix.NewStore(store, types.KeyPrefix(types.PortfolioKeyPrefix))

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
	timestamp := k.GetDateFromBlock(ctx.BlockTime())

	val := k.GetPortfolio(
		ctx,
		req.User,
		req.AssetType,
		timestamp,
	)
	// if !found {
	// 	return nil, status.Error(codes.NotFound, "not found")
	// }

	return &types.QueryGetPortfolioResponse{Portfolio: val}, nil
}
