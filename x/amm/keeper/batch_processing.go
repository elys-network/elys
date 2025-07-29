package keeper

import (
	"cosmossdk.io/store/prefix"
	storetypes "cosmossdk.io/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/elys-network/elys/v7/x/amm/types"
)

func (k Keeper) SetLastSwapRequestIndex(ctx sdk.Context, index uint64) {
	store := ctx.TransientStore(k.transientStoreKey)
	store.Set([]byte(types.TLastSwapRequestIndex), sdk.Uint64ToBigEndian(index))
}

func (k Keeper) GetLastSwapRequestIndex(ctx sdk.Context) uint64 {
	store := ctx.TransientStore(k.transientStoreKey)
	bz := store.Get([]byte(types.TLastSwapRequestIndex))
	if bz == nil {
		return 0
	}
	return sdk.BigEndianToUint64(bz)
}

// SetSwapExactAmountInRequests stores swap exact amount in request
func (k Keeper) SetSwapExactAmountInRequests(ctx sdk.Context, msg *types.MsgSwapExactAmountIn, index uint64) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountInKey))
	b := k.cdc.MustMarshal(msg)
	store.Set(types.TKeyPrefixSwapExactAmountIn(msg, index), b)
}

// DeleteSwapExactAmountInRequest removes a swap exact amount in request
func (k Keeper) DeleteSwapExactAmountInRequest(ctx sdk.Context, msg *types.MsgSwapExactAmountIn, index uint64) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountInKey))
	store.Delete(types.TKeyPrefixSwapExactAmountIn(msg, index))
}

// GetAllSwapExactAmountInRequests returns all SwapExactAmountIn requests
func (k Keeper) GetAllSwapExactAmountInRequests(ctx sdk.Context) (list []types.MsgSwapExactAmountIn) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountInKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgSwapExactAmountIn
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetFirstSwapExactAmountInRequest(ctx sdk.Context, sprefix []byte) (*types.MsgSwapExactAmountIn, uint64) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountInKey))
	iterator := storetypes.KVStorePrefixIterator(store, sprefix)

	defer iterator.Close()
	if !iterator.Valid() {
		return nil, 0
	}
	var val types.MsgSwapExactAmountIn
	k.cdc.MustUnmarshal(iterator.Value(), &val)
	key := iterator.Key()
	return &val, sdk.BigEndianToUint64(key[len(key)-8:])
}

// SetSwapExactAmountOutRequests stores swap exact amount out request
func (k Keeper) SetSwapExactAmountOutRequests(ctx sdk.Context, msg *types.MsgSwapExactAmountOut, index uint64) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountOutKey))
	b := k.cdc.MustMarshal(msg)
	store.Set(types.TKeyPrefixSwapExactAmountOut(msg, index), b)
}

// DeleteSwapExactAmountOutRequest deletes a swap exact amount out request
func (k Keeper) DeleteSwapExactAmountOutRequest(ctx sdk.Context, msg *types.MsgSwapExactAmountOut, index uint64) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountOutKey))
	store.Delete(types.TKeyPrefixSwapExactAmountOut(msg, index))
}

// GetAllSwapExactAmountOutRequests returns all SwapExactAmountOut requests
func (k Keeper) GetAllSwapExactAmountOutRequests(ctx sdk.Context) (list []types.MsgSwapExactAmountOut) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountOutKey))
	iterator := storetypes.KVStorePrefixIterator(store, []byte{})

	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var val types.MsgSwapExactAmountOut
		k.cdc.MustUnmarshal(iterator.Value(), &val)
		list = append(list, val)
	}

	return
}

func (k Keeper) GetFirstSwapExactAmountOutRequest(ctx sdk.Context, sprefix []byte) (*types.MsgSwapExactAmountOut, uint64) {
	store := prefix.NewStore(ctx.TransientStore(k.transientStoreKey), types.KeyPrefix(types.TSwapExactAmountOutKey))
	iterator := storetypes.KVStorePrefixIterator(store, sprefix)

	defer iterator.Close()
	if !iterator.Valid() {
		return nil, 0
	}
	var val types.MsgSwapExactAmountOut
	k.cdc.MustUnmarshal(iterator.Value(), &val)
	key := iterator.Key()
	return &val, sdk.BigEndianToUint64(key[len(key)-8:])
}
