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

func (k Keeper) SubAccounts(goCtx context.Context, req *types.SubAccountsRequest) (*types.SubAccountsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	var list []types.SubAccount

	address, err := sdk.AccAddressFromBech32(req.Address)
	if err != nil {
		return nil, err
	}
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	prefixStore := prefix.NewStore(store, types.GetAddressSubAccountPrefixKey(address))

	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var val types.SubAccount
		if err := k.cdc.Unmarshal(value, &val); err != nil {
			return err
		}

		list = append(list, val)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &types.SubAccountsResponse{SubAccounts: list, Pagination: pageRes}, nil
}
