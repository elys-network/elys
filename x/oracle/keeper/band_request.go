package keeper

import (
	errorsmod "cosmossdk.io/errors"
	"github.com/cosmos/cosmos-sdk/runtime"
	sdk "github.com/cosmos/cosmos-sdk/types"
	gogotypes "github.com/cosmos/gogoproto/types"
	"github.com/elys-network/elys/x/oracle/types"
)

// SetBandPriceResult saves the BandPrice result
func (k Keeper) SetBandPriceResult(ctx sdk.Context, requestID types.OracleRequestID, result types.BandPriceResult) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.BandPriceResultStoreKey(requestID), k.cdc.MustMarshal(&result))
}

// GetBandPriceResult returns the BandPrice by requestId
func (k Keeper) GetBandPriceResult(ctx sdk.Context, id types.OracleRequestID) (types.BandPriceResult, error) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.BandPriceResultStoreKey(id))
	if bz == nil {
		return types.BandPriceResult{}, errorsmod.Wrapf(types.ErrNotAvailable, "Result for request ID %d is not available.", id)
	}
	var result types.BandPriceResult
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}

// GetLastBandRequestId return the id from the last BandPrice request
func (k Keeper) GetLastBandRequestId(ctx sdk.Context) int64 {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.KeyPrefix(types.LastBandRequestIdKey))
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// SetLastBandRequestId saves the id from the last BandPrice request
func (k Keeper) SetLastBandRequestId(ctx sdk.Context, id types.OracleRequestID) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.KeyPrefix(types.LastBandRequestIdKey),
		k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}

// SetBandRequest saves band request waiting for responses
func (k Keeper) SetBandRequest(ctx sdk.Context, requestID types.OracleRequestID, result types.BandPriceCallData) {
	store := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx))
	store.Set(types.BandRequestStoreKey(requestID), k.cdc.MustMarshal(&result))
}

// GetBandRequest returns band request waiting for responses
func (k Keeper) GetBandRequest(ctx sdk.Context, id types.OracleRequestID) (types.BandPriceCallData, error) {
	bz := runtime.KVStoreAdapter(k.storeService.OpenKVStore(ctx)).Get(types.BandRequestStoreKey(id))
	if bz == nil {
		return types.BandPriceCallData{}, errorsmod.Wrapf(types.ErrNotAvailable, "BandPriceCallData for request ID %d is not available.", id)
	}
	var result types.BandPriceCallData
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}
