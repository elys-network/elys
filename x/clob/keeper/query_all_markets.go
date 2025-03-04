package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) AllMarkets(goCtx context.Context, req *types.AllMarketsRequest) (*types.AllMarketsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var markets []types.PerpetualMarket

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.PerpetualMarketPrefix)

	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var market types.PerpetualMarket
		if err := k.cdc.Unmarshal(value, &market); err != nil {
			return err
		}

		markets = append(markets, market)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.AllMarketsResponse{PerpetualMarkets: markets, Pagination: pageRes}, nil
}
