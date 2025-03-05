package keeper

import (
	"context"
	"cosmossdk.io/store/prefix"
	"fmt"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/address"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/elys-network/elys/x/clob/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (k Keeper) OwnerPerpetuals(goCtx context.Context, req *types.OwnerPerpetualsRequest) (*types.OwnerPerpetualsResponse, error) {
	ctx := sdk.UnwrapSDKContext(goCtx)
	owner, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, err
	}

	var list []types.Perpetual

	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	key := append(types.PerpetualOwnerPrefix, address.MustLengthPrefix(owner.Bytes())...)
	prefixStore := prefix.NewStore(store, key)

	pageRes, err := query.Paginate(prefixStore, req.Pagination, func(key []byte, value []byte) error {
		var p types.Perpetual
		if err := k.cdc.Unmarshal(value, &p); err != nil {
			return err
		}

		list = append(list, p)
		return nil
	})
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	fmt.Println(list)

	return &types.OwnerPerpetualsResponse{List: list, Pagination: pageRes}, nil
}
