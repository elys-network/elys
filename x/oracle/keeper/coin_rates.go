package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/elys-network/elys/x/oracle/types"
	gogotypes "github.com/gogo/protobuf/types"
)

// SetCoinRatesResult saves the CoinRates result
func (k Keeper) SetCoinRatesResult(ctx sdk.Context, requestID types.OracleRequestID, result types.CoinRatesResult) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoinRatesResultStoreKey(requestID), k.cdc.MustMarshal(&result))
}

// GetCoinRatesResult returns the CoinRates by requestId
func (k Keeper) GetCoinRatesResult(ctx sdk.Context, id types.OracleRequestID) (types.CoinRatesResult, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.CoinRatesResultStoreKey(id))
	if bz == nil {
		return types.CoinRatesResult{}, sdkerrors.Wrapf(types.ErrSample, "Result for request ID %d is not available.", id)
	}
	var result types.CoinRatesResult
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}

// GetLastBandRequestId return the id from the last CoinRates request
func (k Keeper) GetLastBandRequestId(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastBandRequestIdKey))
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// SetLastBandRequestId saves the id from the last CoinRates request
func (k Keeper) SetLastBandRequestId(ctx sdk.Context, id types.OracleRequestID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastBandRequestIdKey),
		k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}

// SetBandRequest saves band request waiting for responses
func (k Keeper) SetBandRequest(ctx sdk.Context, requestID types.OracleRequestID, result types.CoinRatesCallData) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.CoinRatesResultStoreKey(requestID), k.cdc.MustMarshal(&result))
}

// SetBandRequest returns band request waiting for responses
func (k Keeper) GetBandRequest(ctx sdk.Context, id types.OracleRequestID) (types.CoinRatesCallData, error) {
	bz := ctx.KVStore(k.storeKey).Get(types.CoinRatesResultStoreKey(id))
	if bz == nil {
		return types.CoinRatesCallData{}, sdkerrors.Wrapf(types.ErrSample, "CoinRatesCallData for request ID %d is not available.", id)
	}
	var result types.CoinRatesCallData
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}
