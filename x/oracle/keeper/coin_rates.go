package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	gogotypes "github.com/gogo/protobuf/types"
	"github.com/elys-network/elys/x/oracle/types"
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
		return types.CoinRatesResult{}, sdkerrors.Wrapf(types.ErrSample,
			"GetResult: Result for request ID %d is not available.", id,
		)
	}
	var result types.CoinRatesResult
	k.cdc.MustUnmarshal(bz, &result)
	return result, nil
}

// GetLastCoinRatesID return the id from the last CoinRates request
func (k Keeper) GetLastCoinRatesID(ctx sdk.Context) int64 {
	bz := ctx.KVStore(k.storeKey).Get(types.KeyPrefix(types.LastCoinRatesIDKey))
	intV := gogotypes.Int64Value{}
	k.cdc.MustUnmarshalLengthPrefixed(bz, &intV)
	return intV.GetValue()
}

// SetLastCoinRatesID saves the id from the last CoinRates request
func (k Keeper) SetLastCoinRatesID(ctx sdk.Context, id types.OracleRequestID) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.KeyPrefix(types.LastCoinRatesIDKey),
		k.cdc.MustMarshalLengthPrefixed(&gogotypes.Int64Value{Value: int64(id)}))
}
